import AjaxService from 'ember-ajax/services/ajax';
import $ from 'jquery';

import ENV from 'front-ember/config/environment';

$.ajaxSetup({    
    xhrFields: {
       withCredentials: true
    },
    crossDomain: true
});

export default AjaxService.extend(ENV.backend);
// export default AjaxService.extend({
//   session: Ember.inject.service(),
//   headers: Ember.computed('session.authToken', {
//     get() {
//       let headers = {};
//       const authToken = this.get('session.authToken');
//       if (authToken) {
//         headers['auth-token'] = authToken;
//       }
//       return headers;
//     }
//   })
// });
