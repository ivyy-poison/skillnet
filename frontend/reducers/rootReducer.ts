// import { combineReducers } from 'redux';
// import userReducer from './userReducer';

// const rootReducer = combineReducers({
//   user: userReducer,
//   // Add more reducers here if needed
// });

// export default rootReducer;

import { combineReducers, Reducer } from 'redux';
import { userReducer } from './userReducer';
import { UserState } from '../types';

type RootState = {
    user: UserState;
    // Add more state types here if needed
};

const rootReducer: Reducer<RootState> = combineReducers<RootState>({
    user: userReducer,
    // Add more reducers here if needed
});
export {RootState, rootReducer};

