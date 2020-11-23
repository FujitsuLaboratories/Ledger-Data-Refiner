/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

$(function () {
  $("input").val("")
  $("#select-paging").val(10)
  var historyParam, historyKeyName;

  historyParam = location.search
  historyKeyName = historyParam.slice((historyParam.indexOf("=") + 1), historyParam.length)

  if (historyKeyName !== "") {
    $("#tab-click li").eq(1).addClass("blocks-active").siblings().removeClass("blocks-active")
  } else {
    $("#tab-click li").eq(0).addClass("blocks-active").siblings().removeClass("blocks-active")
  }

  // get channelid
  $.ajax({
    url: channelIdUrl,
    timeout: 3000
  })
    .done(function (data) {
      var getData, channelID, createrParam;
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
      createrParam = { "channelId": channelID }

      //chaincode name data
      nameSelectData($("#chaincode-name-dockey"), chaincodeNameUrl + channelID)
      //creator of tx data
      nameSelectData($("#transaction-creater"), transactionCreateUrl + channelID)

      //check if url parameter exists
      if (historyKeyName !== "") {

        var $tableThtrFirst = $("#block-table thead tr th:eq(0)");
        $("#block-number-input").val(historyKeyName)
        $("#select-paging").val(10)
        $tableThtrFirst.removeClass("sorting_desc sorting_asc").addClass("sorting_desc")

        $("#table-search>.form-inline").eq(1).removeClass("display-none").siblings().addClass("display-none")
        $("#block-table tbody tr").remove()
        $("#paging-length, #entries-info, #entries-page").removeClass("display-none")

        initTableData(historyKeyUrl, channelID, "historyKey")

      } else {

        //initialize the table of Time Range
        initTableData(historyTimeUrl, channelID, "historyTime")

      }

      //  select the data, show the number of informations on each page.
      $("#select-paging").on("change", function () {

        var $tabTxt = $("#tab-click .blocks-active").index(),
          postUrl,
          choiceTxt;

        if ($tabTxt == 0) {           // Time Range Part
          postUrl = historyTimeUrl;
          choiceTxt = "historyTime";
        } else if ($tabTxt == 1) {     // Doc Key Part
          postUrl = historyKeyUrl;
          choiceTxt = "historyKey";
        }

        initTableData(postUrl, channelID, choiceTxt)

      })

      // sort by Block Height
      $("#block-table thead tr th:eq(0)").on("click", debounce(function () {
        ($(this).hasClass("sorting_desc")) ? $(this).removeClass("sorting_desc").addClass("sorting_asc") : $(this).removeClass("sorting_asc").addClass("sorting_desc")

        var postUrl,
          choiceTxt,
          $tabTxt;
        $tabTxt = $("#tab-click .blocks-active").index()

        if ($tabTxt == 0) {  // Time Range Part
          postUrl = historyTimeUrl
          choiceTxt = "historyTime"
        } else if ($tabTxt == 1) {  // Doc Key Part
          postUrl = historyKeyUrl
          choiceTxt = "historyKey"
        }

        initTableData(postUrl, channelID, choiceTxt)

      }))

      // search by Time Range
      $("#block-time-search").on("click", debounce(function () {
        initTableData(historyTimeUrl, channelID, "historyTime")
      }))

      // search by history doc key
      $("#block-number-search").off("click").on("click", debounce(function () {
        var $input = $("#block-number-input").val().trim()
        if ($input == "") {
          alert("The input can't be empty!")
          return false;
        }
        $("#block-table thead tr th:eq(0)").addClass("sorting_desc")
        $("#paging-length, #entries-info, #entries-page").removeClass("display-none")
        initTableData(historyKeyUrl, channelID, "historyKey")

      }))

      // click the word 'Doc Key' in the table, show some details
      $("#block-table").off("click").on("click", ".history-key", function () {

        var $elem = $(this).parent().parent().children();
        var options = {
          collapsed: false,
          withQuotes: true
        };

        var $form = $("#block-form .block-detail")

        for (var i = 0; i < $form.length; i++) {

          if (i !== 8) {
            $form.eq(i).children().eq(0).text($elem.eq(i).text())
          } else {
            $("#doc-content").jsonViewer(JSON.parse($elem.eq(i).children().text()), options);
          }

        }


      })

      // click on the icon to implement the copy function
      copy($("#block-form"), ".fa-copy")
      copy2($("#history-block-table"), ".history-content")

      // click the tab element, select the corresponding content part.
      $("#tab-click").find("li").on("click", function () {
        var $thisIndex = $(this).index(),
          $pageElement = $("#paging-length,#entries-info,#entries-page"),
          $firstTableTrTh = $("#block-table thead tr th:eq(0)");

        $("#block-load").remove()
        $("#table-search").addClass("animated flipInX")
        setTimeout(function () { $("#table-search").removeClass("animated flipInX") }, 1000)
        $("#table-search>.form-inline").eq($thisIndex).removeClass("display-none").siblings().addClass("display-none")
        $(this).addClass("blocks-active").siblings().removeClass("blocks-active")
        $("#select-paging").val(10)
        $("#block-table tbody tr").remove()
        $firstTableTrTh.removeClass("sorting_desc sorting_asc").addClass("sorting_desc")

        if ($thisIndex == 0) { // time range part
          $("#start-time,#end-time").val("")
          $pageElement.removeClass("display-none")
          initTableData(historyTimeUrl, channelID, "historyTime")

        } else if ($thisIndex == 1) { // Doc Key part
          $firstTableTrTh.removeClass("sorting_desc sorting_asc")
          $("#transaction-creater,#chaincode-name-dockey").val("Please select")

          $("#block-number-input").val("")
          $pageElement.addClass("display-none")
          $("#block-table").append('<tr><td class="text-center" colspan="11">Input Doc Key to Query Its History Please...</td></tr>')

        }

      })

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      alert("Request failed: " + textStatus);
    });


})


/*
* Initialize the table content.
* The table of Time and Doc key, table codes.
* use AJAX to get the data.
*/

function initTableData(url, channelID, choiceTxt) {
  loading($(".block-table-contain"))

  var postParams = {},
    sizePageParam = Number($("#select-paging").val());


  if (choiceTxt == "historyTime") {
    postParams.channel_id = channelID
    postParams.sort_by = $("#block-table thead tr th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC"
    postParams.start_time = $("#start-time").val().trim()
    postParams.end_time = $("#end-time").val().trim()
    postParams.page_num = 1
    postParams.page_count = sizePageParam


  } else if (choiceTxt == "historyKey") {

    postParams.channel_id = channelID
    postParams.sort_by = $("#block-table thead tr th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC"
    postParams.doc_key = $("#block-number-input").val().trim()
    postParams.create_msp = ($("#transaction-creater").val() == "Please select") ? "" : $("#transaction-creater").val()
    postParams.chaincode_name = $("#chaincode-name-dockey").val() == "Please select" ? "" : $("#chaincode-name-dockey").val()
    postParams.page_num = 1
    postParams.page_count = sizePageParam

  }

  $.ajax({
    url: url,
    method: "POST",
    contentType: 'application/json',
    dataType: "json",
    data: JSON.stringify(postParams),
  })
    .done(function (data) {
      $("#block-load").remove();
      $("#block-table tbody tr").remove()
      // table init
      var getData,
        tableLen,
        tableData;
      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }

      if (getData.code !== 200) {
        $("#block-table").append('<tr><td class="text-center" colspan="11">' + getData.msg + '</td></tr>')
        return false;
      }

      tableLen = getData.data.count;
      tableData = getData.data.history;

      $("#block-table").createTable({ tableData: tableData, tableLength: tableLen, tableId: "#block-table", pageId: "#table-paginate", tableInfoId: "#table-info", sizePage: sizePageParam, currentPage: 1, channelId: channelID, choiceContent: choiceTxt, previewUrl: url })

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      // alert( "Request failed: " + textStatus );
      $("#block-load .loading").empty()
      $("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')
    });

}


/*
* get the data of Chaincode Name and Creator of Tx
* @param {Object} $elem  element Id
* @param {String} key  for getting data
*/
function nameSelectData($elem, url) {
  $.ajax({
    url: url,
  })
    .done(function (data) {
      var getData,
        arrayData,
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
      arrayData = getData.data
      getDataLen = arrayData.length

      $elem.find("option:not(:first-child)").remove()

      for (var i = 0; i < getDataLen; i++) {
        strHtml += '<option value="' + arrayData[i] + '">' + arrayData[i] + '</option>';
      }

      $elem.append(strHtml)

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      alert("Request failed: " + textStatus);
    });
}
