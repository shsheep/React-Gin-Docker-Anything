import React, { Component } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import '../index.css';

const hostAddress = "http://localhost:8000"

class SLU extends Component {
    constructor(props) {
        super(props);
        this.state = {
            disabled: true,
            file: null,
            result: null,
            command: null,
        };
    }

    updateFile(e) {
        let file = e.target.files[0]
        this.setState({file: file})
    }

    updateCommand(value) {
        this.setState({
            command: value,
          });
    }

    updateResult(value) {
        this.setState({
            result: value,
        })
        console.log("I'm from updateResult")
        console.log(this.state.result)
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
        console.log(file_to_test)
        console.log(formData)
        console.log(config)
        if (this.state.file === '') {
            alert("파일을 올려주세요")
            return;
        }
        axios.post(hostAddress + "/api/slu-upload", formData, config)
            .then((res) => {
                console.log("Whatever Success")
                console.log(res)
            })
    }

    parseResult() {
        console.log("parseResult activated")
        axios.get(hostAddress + "/api/slu-get-result").then((data) => {
            const tmp = data
            console.log(tmp)
            this.updateResult(data.data.result)
        })
    }

    async handleUploadTest() {
        if (this.state.command === '') {
          alert("명령어가 비었습니다!");
          return;
        }

        await axios.post(hostAddress + "/api/slu-write-get-result/", {
          content: this.state.command,
        });
      }

    render() {
        return (
            <div className="container">
                <hr/>
                <div className="form-group">
                    <input align="middle" type="file" name="file" className="form-control-file" id="exampleInputFile" aria-describedby="fileHelp" onChange={(e) => this.updateFile(e)} />
                    <small id="fileHelp" className="form-text text-muted">작성된 리니지M 명령어 파일을 업로드하세요.</small>
                    <button type="button" className="btn btn-success" onClick={(e) => this.handleTest(e)} > 테스트 </button>
                    <button type="button" className="btn btn-success" onClick={(e) => this.parseResult(e)} > 결과보기 </button>
                </div>

                <div className="row-5">
                    <div className="card border-primary">
                        <div className="card-body text-left">
                            <div class="form-group">
                                <textarea
                                    className="form-control"
                                    id="commandTextarea"
                                    rows="5"
                                    placeholder="이곳에 명령어를 직접 입력하세요."
                                    onChange={(e) => { this.updateCommand(e.target.value) }}
                                />
                                <p></p>
                                <button type="button" className="btn btn-success" onClick={(e) => this.handleUploadTest()} > 업로드+테스트 </button>
                            </div>

                            <div className="form-group">
                                <textarea
                                    className="form-control"
                                    disabled={this.state.disabled}
                                    value={this.state.result}
                                    rows="15"
                                    placeholder="S L U  작 동 중 .  .  ."
                                />
                            </div>
                        </div>
                    </div>
                </div>

                <hr />
                <Link to="/">
                    <button className="btn btn-primary btn space">
                        돌아가기
                        </button>
                </Link>
            </div>
        )
    }
}

export default SLU;