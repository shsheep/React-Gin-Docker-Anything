import React, { Component } from 'react'
import {Route} from 'react-router-dom'
import NavBar from './NavBar/NavBar'
import Modules from './Modules/Modules'
import MSD from './Modules/MSD'
import SLU from './Modules/SLU';

class App extends Component {
  render() {
    return (
      <div>
        <NavBar/>
        <Route exact path='/' component={Modules}/>
        <Route exact path='/msd' component={MSD}/>
        <Route exact path='/slu' component={SLU}/>
      </div>
    );
  }
}

export default App;