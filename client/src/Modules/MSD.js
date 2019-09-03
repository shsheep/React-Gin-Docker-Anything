import React, { Component } from 'react';
import axios from 'axios';
import '../index.css';
import { Link } from 'react-router-dom';

const hostAddress = "http://localhost:8000"

class MSD extends Component {
    constructor(props) {
        super(props);
        this.state = {
            disabled: true,
            file: null,
            result: null,
            filename: null,
        }
    }

    updateFile(e) {
        let file = e.target.files[0]
        this.setState({file: file})
    }

    updateResult(value) {
        this.setState({
            result: value,
        })
    }

    updateFilename(value) {
        this.setState({
            filename: value,
        })
    }

    handleTest(e) {
        console.log(this.state, "THE STATE ---- $$$$")
        let file_to_test = this.state.file
        let formData = new FormData()
        formData.append('file', file_to_test)
        const config = {
            headers: {
                'content-type': 'multipart/form-data'
            }
        }

        axios.post(hostAddress + "/api/msd-upload", formData, config)
            .then((data) => {
                const tmp = data
                console.log(tmp)
                this.updateFilename(data.data.filename)
            })
    }

    parseResult() {
        axios.get(hostAddress + "/api/msd-get-result").then((data) => {
            const tmp = data
            console.log(tmp)
            this.updateResult(data.data.result)
        })
    }

    render() {
        return (
            <div className="container">
                <hr />
                <div className="form-group">
                    <input align="middle" type="file" name="file" onChange={(e) => this.updateFile(e)} />
                    <button type="button" className="btn btn-success" onClick={(e) => this.handleTest(e)} > 테스트 </button>
                    <button type="button" className="btn btn-success" onClick={(e) => this.parseResult(e)} > 결과보기 </button>
                </div>

                <div className="row-5">
                    <div className="card border-primary">
                        <div className="card-body text-left">
                            <div className="form-group">
                                <textarea
                                    value={this.state.result}
                                    disabled={this.state.disabled}
                                    rows="5"
                                    className="form-control"
                                    placeholder="M S D 작 동 중 .  .  ."
                                />
                            </div>
                        </div>
                    </div>
                </div>

                <Link to="/">
                    <button className="btn btn-primary btn space">
                        돌아가기
                    </button>
                </Link>
            </div>
        )
    }
}

export default MSD;