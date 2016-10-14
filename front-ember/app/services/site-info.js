import Ember from 'ember';

export default Ember.Service.extend({
  item: null,
  ajax: Ember.inject.service(),
  init() {
    this._super(...arguments);
    this.get('ajax').request('/siteInfo').then(function(rst) {
      document.title = rst.title;      
      this.set('item', rst);
    }.bind(this));
  },
});
