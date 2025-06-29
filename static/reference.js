/**
 * AdTech Cross-Domain Tracking SDK - Optimized Version
 * Version: 1.4.0
 *
 * Features:
 * - Link decoration parameter extraction
 * - Multi-layer identity persistence (cookies, localStorage, fingerprint)
 * - Safari ITP handling
 * - Event deduplication
 * - Cross-domain tracking with parameter propagation
 * - Auto page tracking
 * - SendBeacon support for reliable delivery
 * - ThumbmarkJS integration for device fingerprinting
 *
 * IMPORTANT: By default, tracking parameters are preserved in URLs
 * for seamless cross-domain tracking. Set preserveTrackingParams: false
 * if you want the old behavior of cleaning URLs.
 *
 * Dependencies:
 * - ThumbmarkJS (loaded dynamically): https://github.com/thumbmarkjs/thumbmarkjs
 *
 * Browser Support:
 * - SendBeacon: Chrome 39+, Firefox 31+, Safari 11.1+, Edge 14+
 * - ThumbmarkJS: All modern browsers
 * - Fallback to simple fingerprinting for older browsers
 */

(function (window, document) {
  'use strict';

  // Configuration defaults
  const CONFIG = {
    endpoint: 'https://track.adtech-platform.com',
    cookieName: '_atk_id',
    cookieDomain: null,
    cookieMaxAge: 604800,
    storageKey: '_atk_identity',
    attributionWindow: 604800000,
    sessionTimeout: 1800000,
    enableFingerprint: true,
    enableITPHandling: true,
    enableAutoTracking: true,
    enableDeduplication: true,
    enableCrossDomainPropagation: true,
    preserveTrackingParams: true,
    propagateToDomains: [], // Empty = all domains, or specify: ['checkout.com', 'payments.com']
    deduplicationWindow: 3600000,
    refreshInterval: 82800000,
    thumbmarkUrl:
      'https://cdn.jsdelivr.net/npm/@thumbmarkjs/thumbmarkjs@latest/dist/thumbmark.umd.js',
    debug: false,
  };

  // Tracking parameters
  const PARAMS = {
    SESSION_ID: 'atk_sid',
    CAMPAIGN_ID: 'atk_cid',
    TIMESTAMP: 'atk_ts',
    SIGNATURE: 'atk_sig',
  };

  // Utility functions
  const utils = {
    log: function (...args) {
      if (CONFIG.debug) console.log('[AdTech]', ...args);
    },

    getCookie: function (name) {
      const match = document.cookie.match(new RegExp(`(^| )${name}=([^;]+)`));
      return match ? match[2] : null;
    },

    setCookie: function (name, value, domain, maxAge) {
      const parts = [
        `${name}=${value}`,
        `max-age=${maxAge}`,
        'path=/',
        domain ? `domain=${domain}` : '',
        'SameSite=Lax',
        window.location.protocol === 'https:' ? 'Secure' : '',
      ].filter(Boolean);

      document.cookie = parts.join('; ');
    },

    generateId: function () {
      return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(
        /[xy]/g,
        function (c) {
          const r = (Math.random() * 16) | 0;
          const v = c === 'x' ? r : (r & 0x3) | 0x8;
          return v.toString(16);
        }
      );
    },

    hashString: function (str) {
      let hash = 0;
      for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = (hash << 5) - hash + char;
        hash = hash & hash;
      }
      return Math.abs(hash).toString(36);
    },

    detectCookieDomain: function () {
      const hostname = window.location.hostname;
      if (hostname === 'localhost' || /^[\d.]+$/.test(hostname)) {
        return hostname;
      }

      const parts = hostname.split('.');
      for (let i = parts.length - 2; i >= 0; i--) {
        const domain = '.' + parts.slice(i).join('.');
        document.cookie = `_atk_test=1; domain=${domain}`;
        if (utils.getCookie('_atk_test')) {
          document.cookie = `_atk_test=; domain=${domain}; max-age=0`;
          return domain;
        }
      }
      return hostname;
    },

    getStorage: function (key) {
      try {
        return JSON.parse(
          localStorage.getItem(key) || sessionStorage.getItem(key) || '{}'
        );
      } catch (e) {
        return {};
      }
    },

    setStorage: function (key, value) {
      const data = JSON.stringify(value);
      try {
        localStorage.setItem(key, data);
      } catch (e) {
        try {
          sessionStorage.setItem(key, data);
        } catch (e2) {}
      }
    },
  };

  // Identity management
  class Identity {
    constructor() {
      this.cookieDomain = CONFIG.cookieDomain || utils.detectCookieDomain();
    }

    get() {
      // Try cookie first
      const cookieValue = utils.getCookie(CONFIG.cookieName);
      if (cookieValue) {
        try {
          const identity = JSON.parse(
            atob(cookieValue + '=='.slice((cookieValue.length % 4) % 2))
          );
          if (this.isValid(identity)) return identity;
        } catch (e) {}
      }

      // Try storage
      const stored = utils.getStorage(CONFIG.storageKey);
      if (stored && stored.sid && this.isValid(stored)) {
        return stored;
      }

      return null;
    }

    set(sid, cid, ts) {
      const identity = {
        sid: sid,
        cid: cid,
        ts: ts || Date.now(),
        created: Date.now(),
      };

      // Set cookie
      const encoded = btoa(JSON.stringify(identity)).replace(/=/g, '');
      utils.setCookie(
        CONFIG.cookieName,
        encoded,
        this.cookieDomain,
        CONFIG.cookieMaxAge
      );

      // Set storage
      utils.setStorage(CONFIG.storageKey, identity);

      return identity;
    }

    isValid(identity) {
      if (!identity || !identity.sid || !identity.cid) return false;

      const age = Date.now() - (identity.ts || 0);
      return age <= CONFIG.attributionWindow;
    }

    refresh() {
      const identity = this.get();
      if (identity) {
        identity.refreshed = Date.now();
        this.set(identity.sid, identity.cid, identity.ts);
      }
    }

    clear() {
      document.cookie = `${CONFIG.cookieName}=; domain=${this.cookieDomain}; max-age=0; path=/`;
      localStorage.removeItem(CONFIG.storageKey);
      sessionStorage.removeItem(CONFIG.storageKey);
    }
  }

  // Event deduplication
  class Deduplication {
    constructor() {
      this.sent = new Map();
      this.critical = new Set();
      this.load();
    }

    shouldSend(name, props) {
      if (!CONFIG.enableDeduplication) return true;

      const key = this.getKey(name, props);
      const now = Date.now();

      // Check critical events
      if (['conversion', 'purchase', 'signup'].includes(name)) {
        const criticalKey =
          props.orderId || props.transactionId || props.userId;
        if (criticalKey && this.critical.has(criticalKey)) {
          return false;
        }
      }

      // Check time window
      const lastSent = this.sent.get(key);
      if (lastSent && now - lastSent < CONFIG.deduplicationWindow) {
        return false;
      }

      return true;
    }

    markSent(name, props) {
      const key = this.getKey(name, props);
      this.sent.set(key, Date.now());

      // Mark critical events
      if (['conversion', 'purchase', 'signup'].includes(name)) {
        const criticalKey =
          props.orderId || props.transactionId || props.userId;
        if (criticalKey) {
          this.critical.add(criticalKey);
        }
      }

      this.save();
    }

    getKey(name, props) {
      const significant = {
        conversion: ['orderId', 'value'],
        purchase: ['orderId', 'transactionId'],
        form_submit: ['formId'],
        page_view: ['path'],
      };

      const keys = significant[name] || Object.keys(props).slice(0, 5);
      const data = { name };

      keys.forEach((k) => {
        if (props[k] !== undefined) data[k] = props[k];
      });

      return utils.hashString(JSON.stringify(data));
    }

    load() {
      try {
        const data = utils.getStorage('_atk_dedup');
        if (data.sent) {
          Object.entries(data.sent).forEach(([k, v]) => {
            if (Date.now() - v < CONFIG.deduplicationWindow) {
              this.sent.set(k, v);
            }
          });
        }
        if (data.critical) {
          this.critical = new Set(data.critical);
        }
      } catch (e) {}
    }

    save() {
      const data = {
        sent: {},
        critical: Array.from(this.critical),
      };

      this.sent.forEach((v, k) => {
        if (Date.now() - v < CONFIG.deduplicationWindow) {
          data.sent[k] = v;
        }
      });

      utils.setStorage('_atk_dedup', data);
    }
  }

  // ThumbmarkJS integration
  class FingerprintManager {
    constructor() {
      this.fingerprint = null;
      this.thumbmarkLoaded = false;
      this.loadThumbmark();
    }

    async loadThumbmark() {
      // Check if ThumbmarkJS is already loaded
      if (typeof ThumbmarkJS !== 'undefined') {
        this.thumbmarkLoaded = true;
        return;
      }

      // Load ThumbmarkJS dynamically
      return new Promise((resolve) => {
        const script = document.createElement('script');
        script.src =
          CONFIG.thumbmarkUrl ||
          'https://cdn.jsdelivr.net/npm/@thumbmarkjs/thumbmarkjs@latest/dist/thumbmark.umd.js';
        script.async = true;
        script.onload = () => {
          this.thumbmarkLoaded = true;
          resolve();
        };
        script.onerror = () => {
          utils.log('Failed to load ThumbmarkJS, fingerprinting disabled');
          resolve(); // Continue without fingerprinting
        };
        document.head.appendChild(script);
      });
    }

    async generate() {
      // Return cached fingerprint if available
      if (this.fingerprint) return this.fingerprint;

      // Ensure ThumbmarkJS is loaded
      if (!this.thumbmarkLoaded) {
        await this.loadThumbmark();
      }

      // Generate fingerprint
      try {
        if (typeof ThumbmarkJS !== 'undefined') {
          this.fingerprint = await ThumbmarkJS.getFingerprint();
          return this.fingerprint;
        }
      } catch (error) {
        utils.log('Fingerprint generation failed:', error);
      }

      // Fallback to simple fingerprint if ThumbmarkJS fails
      return this.generateFallback();
    }

    generateFallback() {
      // Simple fallback fingerprint
      const data = {
        ua: navigator.userAgent,
        lang: navigator.language,
        tz: new Date().getTimezoneOffset(),
        screen: screen.width + 'x' + screen.height,
        platform: navigator.platform,
      };
      return utils.hashString(JSON.stringify(data));
    }
  }

  // Main tracker class
  class AdTechTracker {
    constructor(config) {
      Object.assign(CONFIG, config);

      this.identity = new Identity();
      this.dedup = new Deduplication();
      this.fingerprint = new FingerprintManager();
      this.session = null;
      this.queue = [];
      this.initialized = false;

      // Auto-init
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => this.init());
      } else {
        this.init();
      }
    }

    async init() {
      if (this.initialized) return;

      // Extract parameters
      const params = this.extractParams();

      if (params) {
        // New session from chatbot
        this.session = this.identity.set(params.sid, params.cid, params.ts);
        this.session.isNew = true;

        // Optionally clean URL (disabled by default for cross-domain tracking)
        if (!CONFIG.preserveTrackingParams) {
          this.cleanUrl();
        }
      } else {
        // Returning user
        this.session = this.identity.get();

        if (!this.session) {
          // Create organic session
          const sid = utils.generateId();
          this.session = this.identity.set(sid, 'organic', Date.now());
          this.session.isNew = true;
        }
      }

      // Add fingerprint
      if (CONFIG.enableFingerprint) {
        this.session.fp = await this.fingerprint.generate();
      }

      // Start ITP handling
      if (CONFIG.enableITPHandling && this.isSafari()) {
        this.startITPHandler();
      }

      // Setup auto tracking
      if (CONFIG.enableAutoTracking) {
        this.setupAutoTracking();
      }

      // Setup cross-domain parameter propagation
      if (CONFIG.enableCrossDomainPropagation) {
        this.setupCrossDomainPropagation();
      }

      // Process queue
      this.initialized = true;
      this.processQueue();

      // Track session start
      if (this.session.isNew) {
        this.track('session_start', {
          referrer: document.referrer,
          url: window.location.href,
        });
      }

      utils.log('Initialized', this.session);
    }

    extractParams() {
      const params = new URLSearchParams(window.location.search);
      const sid = params.get(PARAMS.SESSION_ID);
      const cid = params.get(PARAMS.CAMPAIGN_ID);

      if (!sid || !cid) return null;

      const ts = parseInt(params.get(PARAMS.TIMESTAMP)) || Date.now();
      const age = Date.now() - ts;

      if (age > CONFIG.attributionWindow) return null;

      return { sid, cid, ts };
    }

    cleanUrl() {
      const url = new URL(window.location);
      let cleaned = false;

      Object.values(PARAMS).forEach((param) => {
        if (url.searchParams.has(param)) {
          url.searchParams.delete(param);
          cleaned = true;
        }
      });

      if (cleaned) {
        window.history.replaceState({}, '', url.toString());
      }
    }

    setupCrossDomainPropagation() {
      // Get current tracking params
      const urlParams = new URLSearchParams(window.location.search);
      const trackingParams = {};

      // Extract all tracking parameters
      Object.values(PARAMS).forEach((param) => {
        if (urlParams.has(param)) {
          trackingParams[param] = urlParams.get(param);
        }
      });

      // If no tracking params in URL, use current session
      if (Object.keys(trackingParams).length === 0 && this.session) {
        trackingParams[PARAMS.SESSION_ID] = this.session.sid;
        trackingParams[PARAMS.CAMPAIGN_ID] = this.session.cid;
        trackingParams[PARAMS.TIMESTAMP] = this.session.ts;
      }

      // Only proceed if we have tracking params
      if (Object.keys(trackingParams).length === 0) return;

      // Store params for form handling
      window._atkTrackingParams = trackingParams;

      // Add params to all links on click
      document.addEventListener('click', (e) => {
        const link = e.target.closest('a');
        if (!link || !link.href) return;

        try {
          const url = new URL(link.href);

          // Check if we should propagate to this domain
          if (this.shouldPropagateToDomain(url.hostname)) {
            // Add tracking params to the link
            Object.entries(trackingParams).forEach(([key, value]) => {
              if (value) url.searchParams.set(key, value);
            });

            link.href = url.toString();
            utils.log('Added tracking params to link:', url.hostname);
          }
        } catch (e) {
          // Invalid URL, skip
        }
      });

      // Handle form submissions (GET forms)
      document.addEventListener('submit', (e) => {
        const form = e.target;
        if (form.method.toLowerCase() === 'get') {
          try {
            const url = new URL(form.action, window.location.origin);
            if (this.shouldPropagateToDomain(url.hostname)) {
              Object.entries(trackingParams).forEach(([key, value]) => {
                if (value && !form.elements[key]) {
                  const input = document.createElement('input');
                  input.type = 'hidden';
                  input.name = key;
                  input.value = value;
                  form.appendChild(input);
                }
              });
              utils.log('Added tracking params to form:', url.hostname);
            }
          } catch (e) {}
        }
      });
    }

    shouldPropagateToDomain(hostname) {
      // Always propagate to same domain
      if (hostname === window.location.hostname) return true;

      // If specific domains configured, only propagate to those
      if (CONFIG.propagateToDomains && CONFIG.propagateToDomains.length > 0) {
        return CONFIG.propagateToDomains.some(
          (domain) => hostname === domain || hostname.endsWith('.' + domain)
        );
      }

      // By default, propagate to all external domains
      return true;
    }

    track(name, props = {}) {
      // Check deduplication first
      if (!this.dedup.shouldSend(name, props)) {
        utils.log(`Blocked duplicate: ${name}`);
        return;
      }

      const event = {
        name: name,
        properties: props,
        timestamp: Date.now(),
        session: this.session,
      };

      if (!this.initialized) {
        this.queue.push(event);
        return;
      }

      this.dedup.markSent(name, props);
      this.send(event);
    }

    async send(event) {
      event.properties = {
        ...event.properties,
        url: window.location.href,
        title: document.title,
      };

      const url = CONFIG.endpoint + '/event';
      const data = JSON.stringify(event);

      // Use sendBeacon for reliability (especially on page unload)
      if (navigator.sendBeacon && this.shouldUseBeacon(event)) {
        const blob = new Blob([data], { type: 'application/json' });
        const sent = navigator.sendBeacon(url, blob);

        if (sent) {
          utils.log('Sent via beacon:', event.name);
          return;
        }
      }

      // Fallback to fetch
      try {
        const response = await fetch(url, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: data,
          credentials: 'include',
          keepalive: true, // Important for page unload
        });

        if (!response.ok) throw new Error('Failed');

        utils.log('Sent via fetch:', event.name);
      } catch (e) {
        utils.log('Send failed:', e);
        this.retry(event);
      }
    }

    shouldUseBeacon(event) {
      // Use sendBeacon for events that don't need response
      const beaconEvents = [
        'page_view',
        'page_hidden',
        'page_visible',
        'session_start',
        'session_end',
        'navigation',
        'click',
        'scroll_depth',
        'engagement_time',
      ];
      return beaconEvents.includes(event.name);
    }

    retry(event, attempt = 1) {
      if (attempt > 3) return;

      setTimeout(() => {
        // Try sendBeacon first on retry
        if (navigator.sendBeacon && this.shouldUseBeacon(event)) {
          const blob = new Blob([JSON.stringify(event)], {
            type: 'application/json',
          });
          if (navigator.sendBeacon(CONFIG.endpoint + '/event', blob)) {
            utils.log('Retry sent via beacon:', event.name);
            return;
          }
        }

        // Fallback to fetch
        this.send(event).catch(() => this.retry(event, attempt + 1));
      }, Math.min(1000 * Math.pow(2, attempt), 30000));
    }

    processQueue() {
      while (this.queue.length > 0) {
        const event = this.queue.shift();
        // Mark as sent for deduplication
        this.dedup.markSent(event.name, event.properties);
        this.send(event);
      }
    }

    setupAutoTracking() {
      // Page visibility
      document.addEventListener('visibilitychange', () => {
        this.track(document.hidden ? 'page_hidden' : 'page_visible');
      });

      // Session end
      window.addEventListener('beforeunload', () => {
        const event = {
          name: 'session_end',
          properties: {
            duration: Date.now() - this.session.ts,
            url: window.location.href,
          },
          timestamp: Date.now(),
          session: this.session,
        };

        // Use sendBeacon directly for guaranteed delivery
        if (navigator.sendBeacon) {
          const blob = new Blob([JSON.stringify(event)], {
            type: 'application/json',
          });
          navigator.sendBeacon(CONFIG.endpoint + '/event', blob);
        } else {
          // Fallback: try sync XHR (deprecated but works)
          try {
            const xhr = new XMLHttpRequest();
            xhr.open('POST', CONFIG.endpoint + '/event', false); // Sync
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.send(JSON.stringify(event));
          } catch (e) {}
        }
      });

      // SPA tracking
      const originalPushState = history.pushState;
      history.pushState = (...args) => {
        originalPushState.apply(history, args);
        this.track('navigation', { path: args[2] });
      };

      window.addEventListener('popstate', () => {
        this.track('navigation', { path: window.location.pathname });
      });
    }

    isSafari() {
      const ua = navigator.userAgent.toLowerCase();
      return ua.indexOf('safari') > -1 && ua.indexOf('chrome') === -1;
    }

    startITPHandler() {
      // Refresh identity every 23 hours
      setInterval(() => this.identity.refresh(), CONFIG.refreshInterval);

      // Add tracking params to same-site links
      document.addEventListener('click', (e) => {
        const link = e.target.closest('a');
        if (!link || !link.href) return;

        try {
          const url = new URL(link.href);
          if (url.hostname === window.location.hostname) {
            url.searchParams.set('_atk_sid', this.session.sid);
            link.href = url.toString();
          }
        } catch (e) {}
      });
    }
  }

  // Initialize
  window.AdTechTracker = AdTechTracker;
  window.adTech = new AdTechTracker(window.adTechConfig || {});
})(window, document);
