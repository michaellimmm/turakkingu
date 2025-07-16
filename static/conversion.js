(function (window, document) {
  'use strict';

  // TODO:
  // - right now, we don't count attribution window. get the config from API and save it to localstorage (expired in 1 hours) and calculate in FE
  // - other option don't need to save attribution window but always send to server and let server to decide if conversion is valid or not
  // design question: should we put accepted domain on settings?

  // Configuration defaults
  const CONFIG = {
    endpoint: 'https://zeals-tracker-api.ngrok.app',
    cookieName: '_zt_id',
    cookieMaxAge: 30 * 24 * 60 * 60, // 30 days
    storageKey: '_zt_identity',
    deduplicationKey: '_zt_dedup',
    sessionTimeout: 1800000,
    deduplicationWindow: 3600000,
    refreshInterval: 82800000,
    propagateToDomains: [],
  };

  const PARAMS = {
    TRACKER_ID: 'ztid',
    TIMESTAMP: 'ztts',
  };

  const utils = {
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
        document.cookie = `_zt_test=1; domain=${domain}`;
        if (utils.getCookie('_zt_test')) {
          document.cookie = `_zt_test=; domain=${domain}; max-age=0`;
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
      this.cookieDomain = utils.detectCookieDomain();
    }

    restorePadding(base64) {
      const remainder = base64.length % 4;
      if (remainder === 0) {
        return base64;
      }
      if (remainder === 1) {
        throw new Error("Invalid Base64 string (length mod 4 == 1).");
      }
      return base64 + "=".repeat(4 - remainder);
    }

    get() {
      // Try cookie first
      const cookieValue = utils.getCookie(CONFIG.cookieName);
      if (cookieValue) {
        try {
          const padded = this.restorePadding(cookieValue);
          const jsonString = atob(padded);

          const identity = JSON.parse(jsonString);
          if (this.isValid(identity)) return identity;
        } catch (e) {
          console.error(e);
        }
      }

      // Try storage
      const stored = utils.getStorage(CONFIG.storageKey);
      if (stored && stored.ztid && this.isValid(stored)) {
        return stored;
      }

      return null;
    }

    set(ztid, ts) {
      const identity = {
        ztid: ztid,
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
      if (!identity || !identity.ztid) return false;
      return true;
    }

    refresh() {
      const identity = this.get();
      if (identity) {
        identity.refreshed = Date.now();
        this.set(identity.ztid, identity.ts);
      }
    }

    clear() {
      document.cookie = `${CONFIG.cookieName}=; domain=${this.cookieDomain}; max-age=0; path=/`;
      localStorage.removeItem(CONFIG.storageKey);
      sessionStorage.removeItem(CONFIG.storageKey);
    }
  }

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
          console.log('Failed to load ThumbmarkJS, fingerprinting disabled');
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
        console.log('Fingerprint generation failed:', error);
      }

      return '';
    }
  }

  // Event deduplication
  class Deduplication {
    constructor() {
      this.sent = new Map();
      this.load();
    }

    shouldSend(event) {
      const key = this.getKey(event);
      const now = Date.now();

      // Check time window
      const lastSent = this.sent.get(key);
      if (lastSent && now - lastSent < CONFIG.deduplicationWindow) {
        return false;
      }

      return true;
    }

    markSent(event) {
      const key = this.getKey(event);
      this.sent.set(key, Date.now());
      this.save();
    }

    getKey(event) {
      // Extract URL components for deduplication
      const url = new URL(event.url || window.location.href);

      const keyData = {
        path: url.pathname,
        ztid: event.session?.ztid,
      };

      // Remove undefined values
      Object.keys(keyData).forEach(
        (key) => keyData[key] === undefined && delete keyData[key]
      );

      return utils.hashString(JSON.stringify(keyData));
    }

    load() {
      try {
        const data = utils.getStorage(CONFIG.deduplicationKey);
        if (data.sent) {
          Object.entries(data.sent).forEach(([k, v]) => {
            if (Date.now() - v < CONFIG.deduplicationWindow) {
              this.sent.set(k, v);
            }
          });
        }
      } catch (e) {}
    }

    save() {
      const data = { sent: {} };

      this.sent.forEach((v, k) => {
        if (Date.now() - v < CONFIG.deduplicationWindow) {
          data.sent[k] = v;
        }
      });

      utils.setStorage(CONFIG.deduplicationKey, data);
    }
  }

  class ZealsTracker {
    constructor() {
      this.identity = new Identity();
      this.fingerprint = new FingerprintManager();
      this.dedup = new Deduplication();
    }

    generateUUID() {
      return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
          (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
      );
    }


    async run() {
      const params = this.extractParams();
      if (params) {
        this.session = this.identity.set(params.ztid, params.ts);
        this.session.isNew = true; // flag if data is not come from storage
      } else {
        this.session = this.identity.get();

        if (!this.session) {
          // if we can't find any data from storage
          var uuid = this.generateUUID();
          this.session = this.identity.set(`web-${uuid}`, Date.now());
          this.session.isNew = true; // flag if data is not come from storage
        }
      }

      // insert fingerprint
      this.session.fp = await this.fingerprint.generate();

      if (this.isSafari()) {
        // refresh cookie
        setInterval(() => this.identity.refresh(), CONFIG.refreshInterval);
      }

      this.setupAutoTracking();

      this.setupCrossDomainPropagation();

      // track session start
      if (this.session.isNew) {
        this.track();
      }
    }

    extractParams() {
      const params = new URLSearchParams(window.location.search);
      const ztid = params.get(PARAMS.TRACKER_ID);

      if (!ztid) return null;

      const ts = parseInt(params.get(PARAMS.TIMESTAMP)) || Date.now();

      return { ztid, ts };
    }

    track() {
      const event = {
        url: window.location.href, // Full URL has everything
        timestamp: Date.now(),
        session: this.session,
      };

      if (!this.dedup.shouldSend(event)) {
        console.log('Blocked duplicate:', event.url);
        return;
      }

      this.dedup.markSent(event);
      this.send(event);
    }

    async send(event) {
      const request = {
        track_id: event.session?.ztid,
        fp: event.session?.fp,
        url: event.url,
        published_at: event.timestamp,
      };

      const url = CONFIG.endpoint + '/v1/tracks/events';
      const data = JSON.stringify(request);

      try {
        const response = await fetch(url, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: data,
          keepalive: true, // Important for page unload
        });

        if (!response.ok) throw new Error('Failed');

        console.log('Sent via fetch');
      } catch (e) {
        console.log('Send failed:', e);
      }
    }

    isSafari() {
      const ua = navigator.userAgent.toLowerCase();
      return ua.indexOf('safari') > -1 && ua.indexOf('chrome') === -1;
    }

    // for SPA tracking
    setupAutoTracking() {
      // Page visibility
      document.addEventListener('visibilitychange', () => {
        // Just track URL when page visibility changes
        this.track();
      });

      // SPA tracking
      const originalPushState = history.pushState;
      history.pushState = (...args) => {
        originalPushState.apply(history, args);
        this.track(); // Just track the new URL
      };

      window.addEventListener('popstate', () => {
        this.track(); // Just track the new URL
      });

      // immediately track
      this.track()
    }

    setupCrossDomainPropagation() {
      const urlParams = new URLSearchParams(window.location.search);
      const trackingParams = {};

      Object.values(PARAMS).forEach((param) => {
        if (urlParams.has(param)) {
          trackingParams[param] = urlParams.get(param);
        }
      });

      if (Object.keys(trackingParams).length === 0 && this.session) {
        trackingParams[PARAMS.TRACKER_ID] = this.session.ztid;
        trackingParams[PARAMS.TIMESTAMP] = this.session.ts;
      }

      if (Object.keys(trackingParams).length === 0) return;

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
            console.log('Added tracking params to link:', url.hostname);
          }
        } catch (e) {
          // Invalid URL, skip
        }
      });

      document.addEventListener('submit', (e) => {
        const form = e.target;

        try {
          const url = new URL(form.action, window.location.origin);

          if (this.shouldPropagateToDomain(url.hostname)) {
            if (form.method.toLowerCase() === 'get') {
              // For GET forms: add as hidden inputs
              Object.entries(trackingParams).forEach(([key, value]) => {
                if (value && !form.elements[key]) {
                  const input = document.createElement('input');
                  input.type = 'hidden';
                  input.name = key;
                  input.value = value;
                  form.appendChild(input);
                }
              });
              console.log('Added tracking params to GET form:', url.hostname);
            } else {
              // For POST forms: append to action URL
              Object.entries(trackingParams).forEach(([key, value]) => {
                if (value) url.searchParams.set(key, value);
              });
              form.action = url.toString();
              console.log(
                'Added tracking params to POST form action:',
                url.hostname
              );
            }
          }
        } catch (e) {
          // Invalid form action, skip
        }
      });
    }

    // TODO: should list all the client domain? for security purpose?
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
  }

  const zealsTracker = new ZealsTracker();
  zealsTracker.run();
  console.log('script is loaded');
})(window, document);
