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
import F2SFunctionCreate from './components/functions/create';
import Authentication from './components/authentication/authentication';
import DeleteF2SFunction from './components/functions/delete';
import F2SImages from './components/images/images';

function App() {
  return (
    <div className="f2s-main-bg">
      <Provider store={store}>
        <ConnectivityCheck >
          <Authentication>
          <Router>
            <NavBar />
            <div className="container">
            <Routes>
              <Route path="/f2sfunctions/create" element={ <F2SFunctionCreate /> } />
              <Route path="/f2sfunctions/:id/invoke" element={ <InvokeFunction /> } />
              <Route path="/f2sfunctions/:id/delete" element={ <DeleteF2SFunction /> } />
              <Route path="/f2sfunctions/:id" element={ <F2SFunctionDetails /> } />
              <Route path="/f2sfunctions" element={ <F2SFunctions /> } />
              <Route path="/images" element={ <F2SImages /> } />

              <Route path="/settings" element={ <Settings /> } />
              <Route path="/" element={ <F2SFunctions /> } />
            </Routes>
            </div>
          </Router>
          </Authentication>
        </ConnectivityCheck>
      </Provider>
    </div>
  );
}

export default App;
