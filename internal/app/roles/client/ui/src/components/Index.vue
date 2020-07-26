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
        fileName: '',
        startDate: '',
        stopDate: '',
        loaded: false,
        values: [],
        labels: [],
        showError: false,
      }
    },
    created () {
      this.requestData()
    },
    methods: {
      resetState () {
        this.loaded = false
        this.showError = false
      },
      reformatDate(rawDate) {
        let dd, d, t;
        dd = new Date(rawDate)
        d = '' + dd.getFullYear() + '-' + ((dd.getMonth() + 1) > 9 ? '' : '0') + (dd.getMonth() + 1) + '-' + (dd.getDate() > 9 ? '' : '0') + dd.getDate()
        t = '' + (dd.getHours() > 9 ? '' : 0) + dd.getHours() + ':' + (dd.getMinutes() > 9 ? '' : 0) + dd.getMinutes()+ ':' + (dd.getSeconds() > 9 ? '' : 0) + dd.getSeconds()

        return d + ' ' + t + ' UTC' + (dd.getTimezoneOffset() > 0 ? '-' : '+') + Math.abs(dd.getTimezoneOffset() / 60)
      },
      requestData () {
        this.resetState()
        axios.get(`/api/v1/timeseries`)
          .then(response => {
            console.log(response.data)
            this.values = response.data.values.map(value => value.value)
            this.labels = response.data.values.map(value => this.reformatDate(value.date))
            this.fileName = response.data.fileName
            this.startDate = this.reformatDate(response.data.startDate)
            this.stopDate = this.reformatDate(response.data.stopDate)
            this.loaded = true
          })
          .catch(err => {
            if (err.response.data.error === null || err.response.data.error === '' || err.response.data.error === undefined) {
              this.errorMessage = err.message
            } else {
              this.errorMessage = err.response.data.error
            }
            this.showError = true
          })
      },
    },
  }
</script>

<template>
  <div class="content">
    <div class="container">
      <div class="error-message" v-if="showError">
        Unexpected error occurred while querying back-end: {{ errorMessage }}
      </div>
      <hr>
      <h1 class="title" v-if="loaded">{{ fileName }}</h1>
      <div class="Chart__container" v-if="loaded">
        <div class="Chart__title">
          Chart built based on values within the following period: <span>{{ startDate }} - {{ stopDate }}</span>
          <hr>
        </div>
        <div class="Chart__content">
          <line-chart v-if="loaded" :chart-data="values" :chart-labels="labels"></line-chart>
        </div>
      </div>
    </div>
  </div>
</template>
