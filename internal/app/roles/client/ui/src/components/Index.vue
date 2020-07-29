<script>
  import axios from 'axios'
  import LineChart from '@/components/LineChart'
  export default {
    name: "IndexPage",
    components: {
      LineChart
    },
    props: {},
    data () {
      return {
        activeFileName: '',
        files: [],
        labels: [],
        valuesLoaded: false,
        loading: false,
        periodStart: '',
        periodStop: '',
        showError: false,
        startDate: '',
        stopDate: '',
        loaded: false,
        values: [],
        labels: [],
        showError: false,
      }
    },
    created () {
      this.listFiles()
        .then(response => {
          console.log('successfully loaded list of files, make initial data request with first file name')
          if (this.files.length > 0) {
            this.activeFileName = this.files[0].fileName
          }
        })
        .catch(err => {
          console.log('failed to load list of files, fallback to request data with empty filename')
        })
    },
    methods: {
      activateFile(f) {
        if (this.activeFileName === f.fileName && this.valuesLoaded) return
        this.activeFileName = f.fileName
        this.periodStart = ''
        this.periodStop = ''
        this.requestData()
      },
      filter() {
        let filter = {
          "fileName": this.activeFileName,
        }
        if (this.periodStart !== '') {
          filter.startDate = this.unixDate(this.periodStart)
        }
        if (this.periodStop !== '') {
          filter.stopDate = this.unixDate(this.periodStop)
        }
        return filter
      },
      listFiles() {
        this.loading = true
        return axios.get('/api/v1/files')
          .then(response => {
            console.log('received list of files', response.data)
            this.files = response.data.map(f => this.parseFileData(f))
            this.loading = false
          })
          .catch(err => {
            console.log('failed to receive list of files: ', err.response.data)
            this.processError(err)
          })
      },
      parseFileData(f) {
        return {
          "fileName": f.fileName,
          "startDate": this.reformatDate(f.startDate),
          "stopDate": this.reformatDate(f.stopDate),
        }
      },
      processError(err) {
        if (err.response.data.error === null || err.response.data.error === '' || err.response.data.error === undefined) {
          this.errorMessage = err.message
        } else {
          this.errorMessage = err.response.data.error
        }
        this.showError = true
        this.loading = false
      },
      resetState() {
        this.valuesLoaded = false
        this.loading = true
        this.showError = false
      },
      reformatDate(rawDate) {
        let dd, d, t;
        dd = new Date(rawDate)
        d = '' + dd.getFullYear() + '-' + ((dd.getMonth() + 1) > 9 ? '' : '0') + (dd.getMonth() + 1) + '-' + (dd.getDate() > 9 ? '' : '0') + dd.getDate()
        t = '' + (dd.getHours() > 9 ? '' : 0) + dd.getHours() + ':' + (dd.getMinutes() > 9 ? '' : 0) + dd.getMinutes()+ ':' + (dd.getSeconds() > 9 ? '' : 0) + dd.getSeconds()

        return d + ' ' + t + ' UTC' + (dd.getTimezoneOffset() > 0 ? '-' : '+') + Math.abs(dd.getTimezoneOffset() / 60)
      },
      requestData (force) {
        if (force === true) {
          this.periodStart = ""
          this.periodStop = ""
        }
        this.resetState()
        axios.get('/api/v1/timeseries', {params: this.filter()})
          .then(response => {
            console.log(response.data)
            this.values = response.data.values.map(value => value.value)
            this.labels = response.data.values.map(value => this.reformatDate(value.date))
            this.fileName = response.data.fileName
            this.startDate = this.reformatDate(response.data.startDate)
            this.stopDate = this.reformatDate(response.data.stopDate)
            this.loading = false
            this.valuesLoaded = true
          })
          .catch(err => {
            this.processError(err)
          })
      },
      unixDate(rawDate) {
        return Math.floor(new Date(rawDate).getTime() / 1000);
      },
    }
  }
</script>
<style>
  .datasource-selector {
    text-align: left;
    margin-left: 50px;
    margin-top: 50px;
  }
  .file-list {
    margin: 0;
    padding-left: 1.2rem;
  }

  .file-list li {
    position: relative;
    list-style-type: none;
    padding-left: 2.5rem;
    margin-bottom: 0.5rem;
  }

  .file-list li:before {
    content: '';
    display: block;
    position: absolute;
    left: 0;
    top: -2px;
    width: 5px;
    height: 11px;
    border-width: 0 2px 2px 0;
    border-style: solid;
    border-color: #00a8a8;
    transform-origin: bottom left;
    transform: rotate(45deg);
  }
  .files {
    float: left;
  }
  .dates {
    float: left;
    margin-left: 50px;
  }
</style>

<template>
  <div class="content">
    <div class="container">
      <div class="error-message" v-if="showError">
        Unexpected error occurred while querying back-end: {{ errorMessage }}
      </div>
      <hr>
      <h1 class="title" v-if="valuesLoaded">{{ activeFileName }}</h1>
      <div class="Chart__container" v-if="loading">
        loading...
      </div>
      <div class="Chart__container" v-if="valuesLoaded">
        <div class="Chart__title">
          Chart built based on values within the following period: <span>{{ startDate }} - {{ stopDate }}</span>
          <hr>
        </div>
        <div class="Chart__content">
          <line-chart v-if="valuesLoaded" :chart-data="values" :chart-labels="labels"></line-chart>
        </div>
      </div>
      <div class="datasource-selector" v-if="!loading">
        <hr>
        <div class="files">
          <h2>Choose CSV file to query:</h2>
          <div v-if="files.length === 0">
            no files available for query, check back-end sevice
          </div>
          <ul class="file-list" v-if="files.length > 0">
            <li v-for="(file,index) in files" :key="index">
              <a href="#" class="file-list-item-action" @click="activateFile(file)">
                {{file.fileName}}: {{file.startDate}} - {{file.stopDate}}
              </a>
            </li>
          </ul>
        </div>
        <div class="dates">
          <h2>Filter source values by date:</h2>
          <label for="start-date" v-text="'Start Date'"></label>
          <datepicker placeholder="Start Date" v-model="periodStart" id="start-date"></datepicker>
          <label for="stop-date" v-text="'Stop Date'"></label>
          <datepicker placeholder="Stop Date" v-model="periodStop" id="stop-date"></datepicker>
          <button title="Apply filter" @click="requestData()">Apply filter</button>
          <button title="Clear filter" @click="requestData(true)">Clear filter</button>
        </div>
      </div>
    </div>
  </div>
</template>
