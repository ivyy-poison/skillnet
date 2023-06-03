// import { configureStore } from '@reduxjs/toolkit';
// import thunk from 'redux-thunk'; // Optional middleware for handling asynchronous actions
// import rootReducer from './reducers/rootReducer';

// const store = configureStore({
//   reducer: rootReducer,
//   middleware: [thunk],
// });

// export default store;

import { configureStore, getDefaultMiddleware } from '@reduxjs/toolkit';
import thunk, { ThunkMiddleware } from 'redux-thunk';
import { AnyAction } from 'redux';
import { rootReducer, RootState } from './reducers/rootReducer';

// Assuming you have a RootState type defined in rootReducer.ts
// type RootState = ReturnType<typeof rootReducer>;

const middleware = [...getDefaultMiddleware<RootState>(), thunk as ThunkMiddleware<RootState, AnyAction>];
const store = configureStore({
    reducer: rootReducer,
    middleware,
});
export default store;
  

