import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { HashRouter } from 'react-router-dom'
import { ToastProvider } from 'react-toast-notifications';

import { Provider } from 'react-redux'
import { createStore } from 'redux'

const goScorbbleStore = (state = false, logIn) => {
  return state = logIn
};

const store = createStore(goScorbbleStore);

ReactDOM.render(
  <HashRouter>
    <ToastProvider autoDismiss="true" autoDismissTimeout="6000" placement="bottom-right">
      <Provider store={store}>
          <App />
      </Provider>
    </ToastProvider>
  </HashRouter>,
  document.getElementById('root')
);
