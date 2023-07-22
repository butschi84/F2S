import logo from './logo.svg';
import './App.css';
import Provider from 'react-redux/es/components/Provider';
import store from './store';
import {  BrowserRouter as Router,  Routes,  Route} from "react-router-dom";
import NavBar from './modules/navbar/navbar'
import F2SFunctions from './components/functions/list';
import InvokeFunction from './components/functions/invoke';
import F2SFunctionDetails from './components/functions/details';
import Settings from './components/settings/settings';
import ConnectivityCheck from './components/connectivity/connectivity'

function App() {
  return (
    <div className="App">
      <Provider store={store}>
        <ConnectivityCheck >
          <Router>
            <NavBar />
            <div className="container">
            <Routes>
              <Route path="/functions/:id/invoke" element={ <InvokeFunction /> } />
              <Route path="/functions/:id" element={ <F2SFunctionDetails /> } />
              <Route path="/functions" element={ <F2SFunctions /> } />

              <Route path="/settings" element={ <Settings /> } />
              <Route path="/" element={ <F2SFunctions /> } />
            </Routes>
            </div>
          </Router>
        </ConnectivityCheck>
      </Provider>
    </div>
  );
}

export default App;
