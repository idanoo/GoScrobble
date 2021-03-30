import {
    REGISTER_SUCCESS,
    REGISTER_FAIL,
    LOGIN_SUCCESS,
    LOGIN_FAIL,
    LOGOUT,
} from "./types";

  import { toast } from 'react-toastify';
  import jwt from 'jwt-decode'
  import AuthService from "../Services/auth.service";
  import eventBus from "./eventBus";

  export const register = (username, email, password) => (dispatch) => {
    return AuthService.register(username, email, password).then(
      (data) => {
        if (data.message) {
          toast.success('Successfully registered. Please sign in');
          dispatch({
            type: REGISTER_SUCCESS,
          });

          return Promise.resolve();
        }

        toast.error(data.error ? data.error: 'An Unknown Error has occurred')
        dispatch({
          type: REGISTER_FAIL,
        });

        return Promise.reject();
      },
      (error) => {
        const message =
          (error.response &&
            error.response.data &&
            error.response.data.message) ||
          error.message ||
          error.toString();

        toast.error(message);

        dispatch({
          type: REGISTER_FAIL,
        });

        return Promise.reject();
      }
    );
  };

  export const login = (username, password) => (dispatch) => {
    return AuthService.login(username, password).then(
      (data) => {
        if (data.token) {
          toast.success('Login Success');
          let user = jwt(data.token)

          dispatch({
            type: LOGIN_SUCCESS,
            payload: { jwt: data.token, sub: user.sub, exp: user.exp },
          });
          return Promise.resolve();
        }

        toast.error(data.error ? data.error: 'An Unknown Error has occurred')
        dispatch({
          type: LOGIN_FAIL,
        });
        return Promise.reject();
      },
      (error) => {
        const message =
          (error.response &&
            error.response.data &&
            error.response.data.error) ||
          error.message ||
          error.toString();

        toast.error('Error: ' + message)
        dispatch({
          type: LOGIN_FAIL,
        });

        return Promise.reject();
      }
    );
  };

  export const logout = (dispatch) => {
      // Clear local data
      AuthService.logout()

      // window.location.pathname("/")
      window.location.reload()

      // TODO; Clear Redux - ENABLE THIS WHEN I FIGURE OUT HOW 2 DISPATCH
      // dispatch({
      //   type: LOGOUT,
      //   payload: {},
      // });

      // // Issue to all listeners to reload
      eventBus.dispatch(LOGOUT);
  };
