export default function authHeader() {
    const auth = localStorage.getItem('user');

    if (auth && auth.jwt) {
      return { Authorization: 'Bearer ' + auth.jwt };
    } else {
      return {};
    }
}
