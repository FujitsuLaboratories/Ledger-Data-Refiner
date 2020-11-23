/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

$(function () {
  //get channel
  $.ajax({
    url: channelIdUrl,
    timeout: 3000
  }).done(function (data) {
    var getData,
      channelID;

    if (typeof (data) == "string") {
      getData = JSON.parse(data);
    } else if (typeof (data) == "object") {
      getData = data;
    }
    if (getData.code !== 200) {
      alert(getData.msg)
      return false;
    }
    channelID = getData.data[0].id

    blockData(blockTimeUrl, channelID)
    transactionData(transactionTimeUrl, channelID)

    setInterval(function () {
      blockData(blockTimeUrl, channelID)
      transactionData(transactionTimeUrl, channelID)
    }, 3 * 1000)

    lineChart(lineChartUrl, channelID, nowDate().slice(0, 10), 0)
    networkTableData(networkUrl + channelID)
    txNum(txStatisticUrl + channelID)

    $("#line-chart-ul li").off("click").on("click", debounce(function () {
      var $index = $(this).index();
      $(this).addClass("active").siblings().removeClass("active");

      if ($index == 0) {
        lineChart(lineChartUrl, channelID, nowDate().slice(0, 10), $index)
      }
      if ($index == 1) {
        lineChart(lineChartUrl2, channelID, nowTo7Date(6).slice(0, 10), $index)
      }
      if ($index == 2) {
        var startTime = nowDate(new Date(new Date(getPreMonthDay(nowDate().slice(0, 10), 1)).getTime() + 24 * 60 * 60 * 1000)).slice(0, 10);
        lineChart(lineChartUrl3, channelID, startTime, $index)
      }
    }))
  })
    .fail(function (jqXHR, textStatus) {
      alert("Request failed: " + textStatus);
    });

  pieChart(pieChartUrl)

})

function networkTableData(url) {
  $.ajax({
    url: url,
  })
    .done(function (data) {
      var getData,
        getDataLen,
        strHtml = "";
      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }

      if (getData.code !== 200) {
        alert(getData.msg)
        return false;
      }

      getData = getData.data
      getDataLen = getData.length

      for (var i = 0; i < getDataLen; i++) {
        strHtml += '<tr><td>' + getData[i].name + '</td><td>' + getData[i].url + '</td><td>' + getData[i].msp + '</td></tr>'
      }

      $("#network-table tbody").empty().append(strHtml)

    })
    .fail(function (jqXHR, textStatus) {
      alert("Request failed: " + textStatus);
    });
}

function txNum(url) {

  $.ajax({
    url: url
  })
    .done(function (data) {
      var getData;
      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }
      if (getData.code !== 200) {
        alert(getData.msg)
        return false;
      }
      getData = getData.data
      $("#block-number").text(getData.block_num)
      $("#transaction-number").text(getData.tx_num)
      $("#node-number").text(getData.node_num)
      $("#contract-number").text(getData.chaincode_num)

    })
    .fail(function (jqXHR, textStatus) {
      alert("Request failed: " + textStatus);
    });

}

// blocks part
function blockData(url, channelID) {
  var postParams = {};
  postParams.channel_id = channelID
  postParams.sort_by = "DESC"
  postParams.start_time = ""
  postParams.end_time = ""
  postParams.page_num = 1
  postParams.page_count = 10

  $.ajax({
    url: url,
    method: "POST",
    contentType: 'application/json',
    dataType: "json",
    data: JSON.stringify(postParams)

  })
    .done(function (data) {
      var getData,
        tableData,
        strHtml = "";

      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }

      if (getData.code !== 200) {
        $("#block-list").empty().append('<li><div class="text-center">' + getData.msg + '</div></li>')
        return false;
      }

      tableData = getData.data.blocks;
      for (var i = 0; i < 5; i++) {
        strHtml += '<li class="block-flex"><div><i class="fa fa-cube"></i></div><div><div><span>Block Hash：</span><span>' + tableData[i].data_hash + '</span></div><div><span>Number of Tx：</span><span>' + tableData[i].tx_count + '</span></div><div><span>Block Height：</span><span>' + tableData[i].block_num + '</span></div><div class="time-style">' + timeDifference(tableData[i].created_at.slice(0, 19).replace(/T/g, " ")) + '</div></div></li>'
      }
      $("#block-list").empty().append(strHtml);

    })
    .fail(function (jqXHR, textStatus) {
      alert("Request failed: " + textStatus);
    });
}

// transaction part
function transactionData(url, channelID) {
  var postParams = {};
  postParams.channel_id = channelID
  postParams.sort_by = "DESC"
  postParams.start_time = ""
  postParams.end_time = ""
  postParams.page_num = 1
  postParams.page_count = 10
  postParams.create_msp = ""

  $.ajax({
    url: url,
    method: "POST",
    contentType: 'application/json',
    dataType: "json",
    data: JSON.stringify(postParams)

  })
    .done(function (data) {
      var getData,
        tableData,
        strHtml = "";

      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }

      if (getData.code !== 200) {
        $("#transaction-list").empty().append('<li><div class="text-center">' + getData.msg + '</div></li>')
        return false;
      }
      tableData = getData.data.txs;
      for (var i = 0; i < 5; i++) {
        strHtml += '<li class="block-flex"><div><i class="fa fa-cube"></i></div><div><div><span>Tx Hash：</span><span>' + tableData[i].tx_hash + '</span></div><div><span>Chaincode：</span><span>' + tableData[i].chaincode_name + '</span></div><div><span>Creator：</span><span>' + tableData[i].create_msp_id + '</span></div><div class="time-style">' + timeDifference(tableData[i].created_at.slice(0, 19).replace(/T/g, " ")) + '</div></div></li>'
      }
      $("#transaction-list").empty().append(strHtml);

    })
    .fail(function (jqXHR, textStatus) {

      alert("Request failed: " + textStatus);
    });
}

// line ajax
function lineChart(url, channelID, startTime, choiceContent) {

  var postParams = {};
  postParams.channel_id = channelID
  if (choiceContent == 0) {
    postParams.start_time = nowDate().slice(0, 10)
  } else {
    postParams.start_time = startTime
    postParams.end_time = nowDate().slice(0, 10)
  }

  $.ajax({
    url: url,
    method: "POST",
    data: JSON.stringify(postParams),
    contentType: 'application/json',
    dataType: "json"
  })
    .done(function (data) {
      var getData,
        xAxisData = [];
      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }

      if (choiceContent == 0) {
        for (var i = 0; i < 24; i++) {
          xAxisData.push(i)
        }
      }
      if (choiceContent == 1) {
        for (var i = 0; i < 7; i++) {
          xAxisData.push(nowTo7Date(i).slice(5, 11))
        }
        xAxisData.reverse()
      }
      if (choiceContent == 2) {
        var n = nowDate(),
          b = getPreMonthDay(nowDate().slice(0, 10), 1),
          c = new Date(n).getTime(),
          d = new Date(b).getTime(),
          len = (c - d) / (24 * 60 * 60 * 1000);
        for (var i = 0; i < len; i++) {
          xAxisData.push(nowTo7Date(i).slice(0, 10))
        }

        xAxisData.reverse()
      }

      var myChart = echarts.init(document.getElementById('dashboard-line-chart'));
      myChart.setOption(lineOption(xAxisData, getData.data));

    })
    .fail(function (jqXHR, textStatus) {
      alert("Request failed: " + textStatus);
    });

}

// line configuration item
function lineOption(xData, seriesData) {

  var option = {
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { show: false },
      data: xData
    },
    yAxis: {

      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { show: false },
    },
    series: [{
      // smooth: {show :true},
      symbol: 'none',
      data: seriesData,
      type: 'line',
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [{
            offset: 0, color: '#275eb0'
          }, {
            offset: 1, color: 'transparent'
          }],
          global: false
        }
      },
      itemStyle: {
        normal: {
          color: "#275eb0",
          lineStyle: {
            color: "#275eb0",
            width: 4
          },
        }
      },

    }]
  };
  return option;
}
// pie ajax
function pieChart(url) {
  $.ajax({
    url: url
  })
    .done(function (data) {
      var getData;
      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }
      var pieData = getData.data,
        pieDataLen = pieData.length,
        legendData = [];
      for (var i = 0; i < pieDataLen; i++) {
        legendData.push(pieData[i].name)
      }

      var myChart = echarts.init(document.getElementById('dashboard-pie-chart'));
      myChart.setOption(pieOption(legendData, pieData));

    })
    .fail(function (jqXHR, textStatus) {
      alert("Request failed: " + textStatus);
    });

}
// pie configuration item
function pieOption(legendData, seriesData) {
  var option = {
    tooltip: {},
    legend: {
      orient: 'horizontal',
      right: 10,
      top: 20,
      bottom: 20,
      icon: 'circle',
      data: legendData,
    },

    series: [{
      type: 'pie',
      radius: '50%',
      data: seriesData
    }],
    color: ['#3520BD', '#275EB0', '#2238C7', '#228FC7', '#20B1BD']
  };
  return option;
}

/*
  caculate time difference
*/

function timeDifference(startTime) {

  var date1 = (new Date()).getTime(),
    date2 = new Date(startTime).getTime() || dayjs(startTime).valueOf(),
    date3 = date1 - date2,
    secondsDate = parseInt(date3 / 1000)

  if (secondsDate < 60) {
    return secondsDate + "seconds ago"
  } else if (secondsDate >= 60) {
    var minDate = parseInt(date3 / (1000 * 60))
    if (minDate < 60) {
      return minDate + "minutes ago"
    } else if (minDate >= 60) {
      var hoursDate = parseInt(date3 / (1000 * 60 * 60))
      if (hoursDate < 24) {
        return hoursDate + "hours ago "
      } else {
        var dayDate = parseInt(date3 / (1000 * 60 * 60 * 24))
        return dayDate + "days ago"
      }

    }

  }

}

function nowDate(val) { // now date
  var date, year, day, nowdatestr, nowMin,
    date = val ? val : new Date();
  year = date.getFullYear();
  month = date.getMonth() + 1 > 9 ? date.getMonth() + 1 : "0" + (date.getMonth() + 1);
  day = date.getDate() > 9 ? date.getDate() : "0" + date.getDate();
  nowMin = ((new Date()).toString()).slice(16, 24);
  nowdatestr = year + "-" + month + "-" + day + " " + nowMin;
  return nowdatestr;
}

function nowTo7Date(i) { // a week ago 
  var date, oneweekdate, year1, month1, day1, dateTime1, nowMin;
  date = new Date();
  nowMin = ((new Date()).toString()).slice(16, 24)
  oneweekdate = new Date(date.getTime() - i * 24 * 3600 * 1000);
  year1 = oneweekdate.getFullYear();
  month1 = oneweekdate.getMonth() + 1 > 9 ? oneweekdate.getMonth() + 1 : "0" + (oneweekdate.getMonth() + 1);
  day1 = oneweekdate.getDate() > 9 ? oneweekdate.getDate() : "0" + oneweekdate.getDate();
  dateTime1 = year1 + "-" + month1 + "-" + day1 + " " + nowMin;

  return dateTime1;
}
//Gets the date of the N months before the current date.
function getPreMonthDay(date, monthNum) {
  var dateArr = date.split('-');
  var year = dateArr[0]; //Gets the year of the current date.
  var month = dateArr[1]; //Gets the month of the current date.
  var day = dateArr[2]; //Gets the day of the current date.
  var days = new Date(year, month, 0);


  days = days.getDate(); //Gets the number of days of the month in the current date.
  var year2 = year;
  var month2 = parseInt(month) - monthNum;
  if (month2 <= 0) {
    year2 = parseInt(year2) - parseInt(month2 / 12 == 0 ? 1 : parseInt(month2) / 12);
    month2 = 12 - (Math.abs(month2) % 12);
  }
  var day2 = day;
  var days2 = new Date(year2, month2, 0);
  days2 = days2.getDate();
  if (day2 > days2) {
    day2 = days2;
  }
  if (month2 < 10) {
    month2 = '0' + month2;
  }

  var nowMin = ((new Date()).toString()).slice(16, 24);

  var t2 = year2 + '-' + month2 + '-' + day2 + ' ' + nowMin;

  return t2;
}
