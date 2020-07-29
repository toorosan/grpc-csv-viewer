<script>
  import { Line } from 'vue-chartjs'

  export default {
    extends: Line,
    props: {
      chartData: {
        type: Array | Object,
        required: false
      },
      chartLabels: {
        type: Array,
        required: true
      }
    },
    data () {
      return {
        options: {
          scales: {
            yAxes: [{
              ticks: {
                beginAtZero: true
              },
              gridLines: {
                display: true
              }
            }],
            xAxes: [ {
              gridLines: {
                display: false
              }
            }]
          },
          legend: {
            display: false
          },
          responsive: true,
          maintainAspectRatio: false
        }
      }
    },
    mounted: function () {
      let bgGradient = this.$refs.canvas.getContext('2d').createLinearGradient(0, 0, 0, 450)
      bgGradient.addColorStop(0, 'rgba(80,4,4,0.75)')
      bgGradient.addColorStop(0.4, 'rgba(156,0,0,0.75)')
      bgGradient.addColorStop(1, 'rgba(255,255,255,0.9)')
      let lineGradient = this.$refs.canvas.getContext('2d').createLinearGradient(0, 0, 0, 450)
      lineGradient.addColorStop(0, 'rgba(80,4,4, 0.8)')
      lineGradient.addColorStop(0.4, 'rgba(156,0,0,0.8)')
      lineGradient.addColorStop(1, 'rgba(255,255,255,1)')
      this.renderChart({
        labels: this.chartLabels,
        datasets: [
          {
            label: 'value',
            borderColor: lineGradient,
            pointBorderColor: lineGradient,
            pointBackgroundColor: lineGradient,
            pointHoverBackgroundColor: lineGradient,
            pointHoverBorderColor: lineGradient,
            pointBorderWidth: 2,
            pointHoverRadius: 2,
            pointHoverBorderWidth: 5,
            pointRadius: 1,
            borderWidth: 1,
            backgroundColor: bgGradient,
            data: this.chartData
          }
        ]
      }, this.options)
    }
  }
</script>
