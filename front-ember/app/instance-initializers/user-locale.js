export function initialize( appInstance ) {
  // appInstance.inject('route', 'foo', 'service:foo');
  var i18n = appInstance.lookup('service:i18n');
  var lang = localStorage.getItem("locale") || navigator.language || navigator.userLanguage || 'en-us';
  i18n.set('locale', lang);
}

export default {
  name: 'user-locale',
  after: "ember-i18n",
  initialize
};
