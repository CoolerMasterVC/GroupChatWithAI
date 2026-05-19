import React from 'react';
import { Provider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import { store, persistor } from './store';
import AppContent from './components/AppContent';
import './App.css';

function App() {
  return (
    <Provider store={store}>
      <PersistGate loading={<div>Загрузка...</div>} persistor={persistor}>
        <AppContent />
      </PersistGate>
    </Provider>
  );
}

export default App;