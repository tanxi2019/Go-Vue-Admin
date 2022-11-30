<template>
  <div class="container">
    <el-row :gutter="32">
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <!--
          @author:风很大
          @description: 柱状图
          @time: 2022/1/20 0020
          -->
          <Echarts :id="id" :options="bar" :className="className"/>
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <!--
          @author:风很大
          @description: 折线图
          @time: 2022/1/20 0020
          -->
          <Echarts :id="idx" :options="line" :className="classx"/>
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <!--
          @author:风很大
          @description:
          @time: 2022/1/20 0020
          -->
          <Echarts :id="idxs" :options="lines" :className="classx"/>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="32">
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <Echarts :id="idb" :options="idbData" :className="className"/>
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <Echarts :id="idb2" :options="idb2Data" :className="className"/>
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <Echarts :id="idb3" :options="idb3Data" :className="className"/>
        </div>
      </el-col>
    </el-row>

  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import Echarts from '@/components/Echarts'

  export default {
    name: 'index',
    computed: {
      ...mapGetters([
        'name'
      ])
    },
    components: { Echarts },
    data() {
      return {
        id: 'bar',
        className: 'bar',
        bar:{
          title: {
            text: '柱状图'
          },
          color :['#2f89cf'], // 柱状图颜色配置
          tooltip: {
            trigger: "axis",
            axisPointer: {
              type: "shadow" // 默认为直线，可选为：'line' | 'shadow'
            }
          },  // 坐标轴指示器状态，坐标轴触发有效
          grid: {
            left: "0%",
            top: "10px",
            right: "0%",
            bottom: "4%",
            containLabel: true
          }, // 调整坐标位置
          xAxis: [
            {
              type: "category",
              data: [
                "旅游行业",
                "教育培训",
                "游戏行业",
                "医疗行业",
                "电商行业",
                "社交行业",
                "金融行业"
              ],
              axisTick: {
                alignWithLabel: true
              },
              axisLabel: {
                textStyle: {
                  color: "rgba(255,255,255,.6)",
                  fontSize: "12" // x坐标轴字体
                }
              },
              axisLine: {
                show: false
              }
            }
          ], // x坐标轴，图例
          yAxis: [
            {
              type: "value",
              axisLabel: {
                textStyle: {
                  color: "rgba(255,255,255,.6)",
                  fontSize: "12"
                }
              },
              axisLine: {
                lineStyle: {
                  color: "rgba(255,255,255,.1)"
                  // width: 1,
                  // type: "solid"
                }
              }, // Y轴线
              splitLine: {
                lineStyle: {
                  color: "rgba(255,255,255,.1)"
                }
              } // X轴线
            }
          ],
          series: [
            {
              name: "直接访问",
              type: "bar",
              barWidth: "45%",
              data: [200, 300, 300, 900, 1500, 1200, 600],
              itemStyle: {
                barBorderRadius: 0, // 柱状图圆角
              }
            }
          ] // 柱状图 柱状高度
        },
        idx: 'line',
        classx: 'line',
        line: {
          title: {
            text: '折线图'
          },
          xAxis: {
            type: 'category',
            data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
          },
          yAxis: {
            type: 'value'
          },
          series: [
            {
              data: [150, 230, 224, 218, 135, 147, 260],
              type: 'line'
            }
          ]
        },
        idxs: 'lines',
        lines: {
          color: ["#00f2f1", "#ed3f35"],
          tooltip: {
            // 通过坐标轴来触发
            trigger: "axis"
          },
          legend: {
            // 距离容器10%
            right: "10%",
            // 修饰图例文字的颜色
            textStyle: {
              color: "#4c9bfd"
            }
            // 如果series 里面设置了name，此时图例组件的data可以省略
            // data: ["邮件营销", "联盟广告"]
          },
          grid: {
            top: "20%",
            left: "3%",
            right: "4%",
            bottom: "3%",
            show: true,
            borderColor: "#012f4a", // 边框颜色
            containLabel: true
          },

          xAxis: {
            type: "category",
            boundaryGap: false,
            data: [
              "1月",
              "2月",
              "3月",
              "4月",
              "5月",
              "6月",
              "7月",
              "8月",
              "9月",
              "10月",
              "11月",
              "12月"
            ],
            // 去除刻度
            axisTick: {
              show: false
            },
            // 修饰刻度标签的颜色
            axisLabel: {
              color: "rgba(255,255,255,.7)"
            },
            // 去除x坐标轴的颜色
            axisLine: {
              show: false
            }
          }, // 图例
          yAxis: {
            type: "value",
            // 去除刻度
            axisTick: {
              show: false
            },
            // 修饰刻度标签的颜色
            axisLabel: {
              color: "rgba(255,255,255,.7)"
            },
            // 修改y轴分割线的颜色
            splitLine: {
              lineStyle: {
                color: "#012f4a"
              }
            }
          },
          series: [
            {
              name: "新增粉丝",
              type: "line",
              stack: "总量",
              // 是否让线条圆滑显示
              smooth: true,
              data:[24, 40, 101, 134, 90, 230, 210, 230, 120, 230, 210, 120],

            },
            {
              name: "新增游客",
              type: "line",
              stack: "总量",
              smooth: true,
              data:[40, 64, 191, 324, 290, 330, 310, 213, 180, 200, 180, 79]

            }
          ]
        },
        idb:'idb',
        idbData:{
          tooltip: {
            trigger: "item",
            formatter: "{a} <br/>{b}: {c} ({d}%)",
            position: function (p) {
              //其中p为当前鼠标的位置
              return [p[0] + 10, p[1] - 10];
            }
          }, // 指示器效果
          legend: {
            top: "90%",
            itemWidth: 10,
            itemHeight: 10,
            data: ["0岁以下", "20-29岁", "30-39岁", "40-49岁", "50岁以上"],
            textStyle: {
              color: "rgba(255,255,255,.5)",
              fontSize: "12"
            }
          },
          series: [
            {
              name: "年龄分布",
              type: "pie",
              center: ["50%", "42%"],
              radius: ["40%", "60%"],
              color: [
                "#065aab",
                "#066eab",
                "#0682ab",
                "#0696ab",
                "#06a0ab",
                "#06b4ab",
                "#06c8ab",
                "#06dcab",
                "#06f0ab"
              ],
              label: {show: false},
              labelLine: {show: false},
              data: [
                {value: 1, name: "0岁以下"},
                {value: 4, name: "20-29岁"},
                {value: 2, name: "30-39岁"},
                {value: 2, name: "40-49岁"},
                {value: 1, name: "50岁以上"}
              ]
            }
          ]
        },
        idb2:'idb2',
        idb2Data:{
          legend: {
            top: "90%",
            itemWidth: 10,
            itemHeight: 10,
            textStyle: {
              color: "rgba(255,255,255,.5)",
              fontSize: "12"
            }
          },
          tooltip: {
            trigger: "item",
            formatter: "{a} <br/>{b} : {c} ({d}%)"
          },
          // 注意颜色写的位置
          color: [
            "#006cff",
            "#60cda0",
            "#ed8884",
            "#ff9f7f",
            "#0096ff",
            "#9fe6b8",
            "#32c5e9",
            "#1d9dff"
          ],
          series: [
            {
              name: "点位统计",
              type: "pie",
              // 如果radius是百分比则必须加引号
              radius: ["10%", "70%"],
              center: ["50%", "42%"],
              roseType: "radius",
              data: [
                {value: 20, name: "云南"},
                {value: 26, name: "北京"},
                {value: 24, name: "山东"},
                {value: 25, name: "河北"},
                {value: 20, name: "江苏"},
                {value: 25, name: "浙江"},
                {value: 30, name: "深圳"},
                {value: 42, name: "广东"}
              ],
              // 修饰饼形图文字相关的样式 label对象
              label: {
                fontSize: 10
              },
              // 修饰引导线样式
              labelLine: {
                // 连接到图形的线长度
                length: 10,
                // 连接到文字的线长度
                length2: 10
              }
            }
          ]
        },
        idb3:'idb3',
        idb3Data:{
          //图标位置
          grid: {
            top: "10%",
            left: "22%",
            bottom: "10%"
          },
          xAxis: {
            show: false
          },
          yAxis: [
            {
              show: true,
              data: ['HTML5', 'CSS3', 'JavaScript', 'VUE', 'NODE'],
              inverse: true,
              axisLine: {
                show: false
              },
              splitLine: {
                show: false
              },
              axisTick: {
                show: false
              },
              axisLabel: {
                color: "#fff",

                rich: {
                  lg: {
                    backgroundColor: "#339911",
                    color: "#fff",
                    borderRadius: 15,
                    // padding: 5,
                    align: "center",
                    width: 15,
                    height: 15
                  }
                }
              }
            },
            {
              show: false,
              inverse: true,
              data: [702, 350, 610, 793, 664],
              axisLabel: {
                textStyle: {
                  fontSize: 12,
                  color: "#fff"
                }
              }
            }
          ],
          series: [
            {
              name: "条",
              type: "bar",
              yAxisIndex: 0,
              data:  [70, 34, 60, 78, 69],
              barCategoryGap: 50,
              barWidth: 14.5,
              itemStyle: {
                normal: {
                  barBorderRadius: 20,
                  color: function (params) {
                    let myColor =  ['#1089E7', '#F57474', '#56D0E3', '#F8B448', '#8B78F6']
                    let num = myColor.length;
                    return myColor[params.dataIndex % num];
                  }
                }
              },
              label: {
                normal: {
                  show: true,
                  position: "inside",
                  formatter: "{c}%"
                }
              }
            },
            {
              name: "框",
              type: "bar",
              yAxisIndex: 1,
              barGap: '-100%',
              data: [100, 100, 100, 100, 100],
              barWidth: 15,
              itemStyle: {
                normal: {
                  color: "none",
                  borderColor: "#00c1de",
                  borderWidth: 1,
                  barBorderRadius: 15
                }
              }
            }
          ]
        }

      }
    }
  }
</script>

<style lang="scss" scoped>
  .bar {
    width: 100%;
    height: 400px;
  }
  .line {
    width: 100%;
    height: 400px;
  }
  .chart-wrapper {
    background: #fff;
    /*padding: 16px 16px 0;*/
    margin-bottom: 25px;
  }


  @media (max-width:1024px) {
    .chart-wrapper {
      padding: 8px;
    }
  }
  .dashboard {
  &-container {
     margin: 30px;
   }
  &-text {
     font-size: 30px;
     line-height: 46px;
   }
  }
</style>
