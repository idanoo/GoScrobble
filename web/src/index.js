import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { HashRouter } from 'react-router-dom';

import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.min.css';

import AuthContextProvider from './Contexts/AuthContextProvider';

ReactDOM.render(
  <AuthContextProvider>
    <HashRouter>
      <ToastContainer
        position="bottom-right"
        autoClose={5000}
        hideProgressBar={false}
        newestOnTop={true}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss={false}
        draggable
        pauseOnHover
        />
        <App />
    </HashRouter>
  </AuthContextProvider>,
  document.getElementById('root')
);
