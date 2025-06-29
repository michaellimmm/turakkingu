(function (window, document) {
  'use strict';

  // Configuration defaults
  const CONFIG = {
    endpoint: 'http://localhost:8080',
    cookieName: '_zt_id',
    cookieMaxAge: 30 * 24 * 60 * 60, // 30 days
    storageKey: '_zt_identity',
    sessionTimeout: 1800000,
    deduplicationWindow: 3600000,
    refreshInterval: 82800000,
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
      this.cookieDomain = utils.detectCookieDomain();
      console.log(this.cookieDomain);
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

  function initTracker() {
    console.log('initialized');
    const indentity = new Identity();
  }

  console.log('script is loaded');
  initTracker();
})(window, document);
