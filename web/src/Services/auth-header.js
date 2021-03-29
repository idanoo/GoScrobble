export default function authHeader() {
    const token = JSON.parse(localStorage.getItem('jwt'));

    if (token) {
      return { Authorization: 'Bearer ' + token };
    } else {
      return {};
    }
}
