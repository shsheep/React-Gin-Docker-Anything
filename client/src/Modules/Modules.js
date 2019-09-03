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
      <Link to="/msd">
        <div className="card text-white bg-secondary mb-3">
          <div className="card-body">
            <h4 className="card-title">Music Speech Detection</h4>
          </div>
        </div>
      </Link>

        <div className="row">
          {this.state.ssmls === null && <p></p>}
          {
            this.state.ssmls && this.state.ssmls.map(ssml => (
              <div key={ssml.id} className="col-sm-12 col-md-4 col-lg-3">
                <Link to={`/${ssml.title}`}>
                  <div className="card text-white bg-success mb-3">
                    <div className="card-body">
                      <h4 className="card-title">{ssml.title}</h4>
                    </div>
                  </div>
                </Link>
              </div>
            ))
          }
        </div>
      </div>
    )
  }
}

export default Modules;