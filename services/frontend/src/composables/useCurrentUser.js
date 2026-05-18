import { ref } from 'vue';
import { getMe } from '../api/users.js';

let loadPromise = null;
let loadedToken = '';
const firstName = ref('');

export function useCurrentUser() {
  return {
    firstName,
    load() {
      const token = localStorage.getItem('token') || '';
      if (!token) {
        firstName.value = '';
        loadPromise = null;
        loadedToken = '';
        return Promise.resolve();
      }

      if (!loadPromise || loadedToken !== token) {
        loadPromise = getMe()
          .then(profile => {
            firstName.value = profile.firstName;
            loadedToken = token;
            return profile;
          })
          .catch((err) => {
            firstName.value = '';
            loadPromise = null;
            loadedToken = '';
            throw err;
          });
      }
      return loadPromise;
    },
    reset() {
      loadPromise = null;
      loadedToken = '';
      firstName.value = '';
    },
  };
}
