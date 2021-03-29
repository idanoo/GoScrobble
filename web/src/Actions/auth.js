import {
    REGISTER_SUCCESS,
    REGISTER_FAIL,
    LOGIN_SUCCESS,
    LOGIN_FAIL,
    SET_MESSAGE,
  } from "./types";
  import { toast } from 'react-toastify'

  import AuthService from "../Services/auth.service";

  export const register = (username, email, password) => (dispatch) => {
    return AuthService.register(username, email, password).then(
      (response) => {
        dispatch({
          type: REGISTER_SUCCESS,
        });

        return Promise.resolve();
      },
      (error) => {
        const message =
          (error.response &&
            error.response.data &&
            error.response.data.message) ||
          error.message ||
          error.toString();

        dispatch({
          type: REGISTER_FAIL,
        });

        dispatch({
          type: SET_MESSAGE,
          payload: message,
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
          dispatch({
            type: LOGIN_SUCCESS,
            payload: { user: data },
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

        // dispatch({
        //   type: SET_MESSAGE,
        //   payload: message,
        // });

        return Promise.reject();
      }
    );
  };

  export const logout = () => () => {
    AuthService.logout();
    window.location.reload();
  };
