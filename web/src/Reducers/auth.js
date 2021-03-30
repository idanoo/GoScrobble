import {
    REGISTER_SUCCESS,
    REGISTER_FAIL,
    LOGIN_SUCCESS,
    LOGIN_FAIL,
    LOGOUT,
  } from "../Actions/types";

  const user = JSON.parse(localStorage.getItem('user'));

  const initialState = user
    ? { isLoggedIn: true, user: user }
    : { isLoggedIn: false, user: null };

  export default function authReducer(state = initialState, action) {
    const { type, payload } = action;

    switch (type) {
      case REGISTER_SUCCESS:
        return {
          ...state,
          isLoggedIn: false,
        };
      case REGISTER_FAIL:
        return {
          ...state,
          isLoggedIn: false,
        };
      case LOGIN_SUCCESS:
        return {
          ...state,
          isLoggedIn: true,
          user: {
            jwt: payload.jwt,
            uuid: payload.sub,
            exp: payload.exp,
          }
        };
      case LOGIN_FAIL:
        return {
          ...state,
          isLoggedIn: false,
        };
      case LOGOUT:
        return {
          ...state,
          isLoggedIn: false,
          user: null,
        };
      default:
        return state;
    }
  }
