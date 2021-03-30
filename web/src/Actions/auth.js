import {
    REGISTER_SUCCESS,
    REGISTER_FAIL,
    LOGIN_SUCCESS,
    LOGIN_FAIL,
  } from "./types";

  import { toast } from 'react-toastify';
  import jwt from 'jwt-decode'
  import AuthService from "../Services/auth.service";

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

  export const logout = () => (dispatch) => {
    AuthService.logout();

    // dispatch({
    //   type: LOGOUT,
    // });

    window.location.reload();
  };
