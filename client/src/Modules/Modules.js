import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import axios from 'axios';
import '../index.css';

class Modules extends Component {
  constructor(props) {
    super(props);

    this.state = {
      modules: null,
    };
  }

  async componentDidMount() {
    // CASE 4
    const modules = (await axios.get('http://localhost:8080/api/ssmls')).data;
    this.setState({
        modules,   
    });
  }

  render() {
    return (
      
      <div className="container">
      <Link to="/media">
        <div className="card text-white bg-secondary mb-3">
          <div className="card-body">
            <h4 className="card-title">Media Analyzer</h4>
          </div>
        </div>
      </Link>

      <Link to="/tab">
        <div className="card text-white bg-secondary mb-3">
          <div className="card-body">
            <h4 className="card-title">Guitar Pro Converter to PDF</h4>
          </div>
        </div>
      </Link>
      </div>
    )
  }
}

export default Modules;