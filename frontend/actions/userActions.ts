// export const LOGIN_REQUEST = 'LOGIN_REQUEST';
// export const LOGIN_SUCCESS = 'LOGIN_SUCCESS';
// export const LOGIN_FAILURE = 'LOGIN_FAILURE';

// export const loginRequest = () => {
//   return { type: LOGIN_REQUEST }
// }

// export const loginSuccess = (user) => {
//   return { type: LOGIN_SUCCESS, payload: user }
// }

// export const loginFailure = (error) => {
//   return { type: LOGIN_FAILURE, payload: error }
// }

// userActions.ts

import { User } from '../types';
export const LOGIN_REQUEST = 'LOGIN_REQUEST';
export const LOGIN_SUCCESS = 'LOGIN_SUCCESS';
export const LOGIN_FAILURE = 'LOGIN_FAILURE';

export type LoginRequest = {
  type: typeof LOGIN_REQUEST;
};

export type LoginSuccess = {
  type: typeof LOGIN_SUCCESS;
  payload: User;
};

export type LoginFailure = {
  type: typeof LOGIN_FAILURE;
  payload: Error;
};

export const loginRequest = (): LoginRequest => ({
  type: LOGIN_REQUEST,
});

export const loginSuccess = (user: User): LoginSuccess => ({
  type: LOGIN_SUCCESS,
  payload: user,
});

export const loginFailure = (error: Error): LoginFailure => ({
  type: LOGIN_FAILURE,
  payload: error,
});


