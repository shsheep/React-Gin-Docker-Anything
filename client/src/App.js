import React, { Component } from 'react'
import {Route} from 'react-router-dom'
import NavBar from './NavBar/NavBar'
import Modules from './Modules/Modules'
import Media from './Modules/Media'
import Tab from './Modules/Tab';

class App extends Component {
  render() {
    return (
      <div>
        <NavBar/>
        <Route exact path='/' component={Modules}/>
        <Route exact path='/media' component={Media}/>
        <Route exact path='/tab' component={Tab}/>
      </div>
    );
  }
}

export default App;