import Ember from 'ember';
import Cookies from 'js-cookie';

export default Ember.Component.extend({
  i18n: Ember.inject.service(),
  siteInfo: Ember.inject.service(),
  locales: Ember.computed('i18n.locale', 'i18n.locales', function() {
    const i18n = this.get('i18n');
    return this.get('i18n.locales').map(function(loc) {
      return {
        id: loc,
        text: i18n.t('language-select.languages.' + loc)
      };
    });
  }),

  actions: {
    setLocale(lang) {
      Cookies.set('locale', lang, { expires: 7 });
      localStorage.setItem("locale", lang);
      this.set('i18n.locale', lang);
    }
  }
});
