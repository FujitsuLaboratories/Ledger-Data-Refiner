/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
$(function () {

  $("input").val("")
  $("#select-paging").val(10)

  // get channelid
  $.ajax({
    url: channelIdUrl,
    timeout: 3000
  })
    .done(function (data) {
      var getData,
        channelID,
        transactionParam,
        transactionBlockNumberName;
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

      // get the data of the Creator of tx
      creatorOftxData(transactionCreateUrl + channelID)

      transactionParam = location.search
      transactionBlockNumberName = transactionParam.slice((transactionParam.indexOf("=") + 1), transactionParam.length)

      //check if url parameter exists
      if (transactionBlockNumberName !== "") {  // has parameters

        var $tableTh = $("#block-table thead tr th");
        $tableTh.eq(0).removeClass("sorting_desc")
        $tableTh.eq(1).addClass("sorting_desc")
        $("#block-number-input").val(transactionBlockNumberName)
        $("#tab-click").find("li").eq(1).addClass("blocks-active").siblings().removeClass("blocks-active")
        $("#table-search>.form-inline").eq(1).removeClass("display-none").siblings().addClass("display-none")

        $("#paging-length, #entries-info, #entries-page").removeClass("display-none")
        initTableData(transactionBlockNumberUrl, channelID, "transactionBlockNumber")

      } else {                                      //no parameters
        //initialize the table of Time Range
        initTableData(transactionTimeUrl, channelID, "transactionTime")
      }

      // click the tab element, select the corresponding content part.
      $("#tab-click").find("li").on("click", debounce(function () {
        var $thisIndex = $(this).index(),
          $tableThtrFirst = $("#block-table thead tr th:eq(0)"),
          $tableThtrSecond = $("#block-table thead tr th:eq(1)"),
          $pageElement = $("#paging-length,#entries-info,#entries-page");

        $("#block-load").remove();
        $("#table-search").addClass("animated flipInX")
        setTimeout(function () { $("#table-search").removeClass("animated flipInX") }, 1000)
        $("#table-search>.form-inline").eq($thisIndex).removeClass("display-none").siblings().addClass("display-none")
        $(this).addClass("blocks-active").siblings().removeClass("blocks-active")

        $("#start-time,#end-time,#block-hash-input").val("")
        $("#select-paging").val(10)
        $("#block-table tbody tr").remove()
        $tableThtrFirst.removeClass("sorting_desc sorting_asc")
        $tableThtrSecond.removeClass("sorting_desc sorting_asc")

        if ($thisIndex == 0) {                        //Time Range part

          $tableThtrFirst.addClass("sorting_desc")
          $pageElement.removeClass("display-none")

          initTableData(transactionTimeUrl, channelID, "transactionTime")

        } else if ($thisIndex == 1) {             //Block Height part

          //   $tableThtrSecond.addClass("sorting_desc")
          $("#block-table").append('<tr><td class="text-center" colspan="8">Input Block Height to Query Tx Please...</td></tr>')
          //no parameters

          $("#block-number-input").val("")
          $pageElement.addClass("display-none")

        } else if ($thisIndex == 2) {                //Tx Hash part

          $("#block-hash-input").val("")
          $pageElement.addClass("display-none")
          $("#block-table").append('<tr><td class="text-center" colspan="8">Input Tx Hash to Query Tx Please...</td></tr>')

        }
      }))

      // search by Time Range
      $("#block-time-search").on("click", debounce(function () {
        initTableData(transactionTimeUrl, channelID, "transactionTime")
      }))

      // search by Blcok Height
      $("#block-number-search").off("click").on("click", debounce(function () {
        var reg = /^\d+$/g
        var $input = $("#block-number-input").val().trim();

        if (!reg.test($input)) {
          alert("Please fill in non-negative integers!")
          return false;
        }
        $("#block-table thead tr th:eq(1)").addClass("sorting_desc")
        $("#paging-length,#entries-info,#entries-page").removeClass("display-none")
        initTableData(transactionBlockNumberUrl, channelID, "transactionBlockNumber")

      }))

      // search by Data Hash
      $("#block-hash-search").off("click").on("click", debounce(function () {
        var $input = $("#block-hash-input").val().trim();
        if ($input == "") {
          alert("Cannot be empty!")
          return false;
        }

        txHashTableData(transactionHashUrl + channelID + "/" + $input)

      }))

      //  select the data, show the number of informations on each page.
      $("#select-paging").on("change", function () {
        var $tabTxt = $("#tab-click .blocks-active").index(),
          postUrl,
          choiceTxt;

        if ($tabTxt == 0) {  // Time Range Part

          postUrl = transactionTimeUrl
          choiceTxt = "transactionTime"

        } else if ($tabTxt == 1) {  // Block Height Part

          postUrl = transactionBlockNumberUrl
          choiceTxt = "transactionBlockNumber"

        }

        initTableData(postUrl, channelID, choiceTxt)

      })

      // the part of Time Range sort by Block Height
      // the part of Block Height sort by Tx Number
      $("#block-table thead tr th:eq(0),#block-table thead tr th:eq(1)").on("click", debounce(function () {

        var $self = $(this);
        ($self.hasClass("sorting_desc")) ? $self.removeClass("sorting_desc").addClass("sorting_asc") : $self.removeClass("sorting_asc").addClass("sorting_desc")

        var postUrl,
          choiceTxt,
          $tabTxt = $("#tab-click .blocks-active").index();

        if ($tabTxt == 0) {  // Time Range part

          postUrl = transactionTimeUrl
          choiceTxt = "transactionTime"

        } else if ($tabTxt == 1) {  // Block Height part

          postUrl = transactionBlockNumberUrl
          choiceTxt = "transactionBlockNumber"

        }

        initTableData(postUrl, channelID, choiceTxt)

      }))

      // click the word 'datahash' in the table, show some details
      $("#block-table").off("click").on("click", ".data-hash", debounce(function () {

        $('#read-set, #write-set, #endorser-signature').empty()

        var transactionParams = {}
        transactionParams.channelId = channelID
        transactionParams.txHash = $(this).text()
        //get data of tx hash
        txHashDetails(transactionHashDetailsUrl + channelID + "/" + $(this).text(), transactionParams)

      }))

      // click on the icon to implement the copy function
      copy($("#block-form2"), ".fa-copy")
    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      alert("Request failed: " + textStatus);
    });



})

/*
* Initialize the table content.
* The table of Time Range , table codes.
* use AJAX to get the data.
*/

function initTableData(url, channelID, choiceTxt) {
  $("#block-load").remove();
  loading($(".block-table-contain"))

  var postParams = {},
    sizePageParam = Number($("#select-paging").val());

  if (choiceTxt == "transactionBlockNumber") {
    postParams.channel_id = channelID
    postParams.page_num = 1
    postParams.page_count = sizePageParam
    postParams.block_num = Number($("#block-number-input").val().trim())
    postParams.sort_by = $("#block-table tr th:eq(1)").hasClass("sorting_desc") ? "DESC" : "ASC"

  } else if (choiceTxt == "transactionTime") {

    postParams.channel_id = channelID
    postParams.sort_by = $("#block-table tr th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC"
    postParams.start_time = $("#start-time").val().trim()
    postParams.end_time = $("#end-time").val().trim()
    postParams.page_num = 1
    postParams.page_count = sizePageParam
    postParams.create_msp = $("#transaction-creater").val() == "Please select" ? "" : $("#transaction-creater").val()

  }

  $.ajax({
    url: url,
    method: "POST",
    data: JSON.stringify(postParams),
    contentType: 'application/json',
    dataType: "json",
  })
    .done(function (data) {
      $("#block-load").remove();
      $("#block-table tbody tr").remove()
      // table 
      var getData,
        tableLen,
        tableData;
      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }
      if (getData.code !== 200) {
        $("#block-table").append('<tr><td class="text-center" colspan="8">' + getData.msg + '</td></tr>')
        return false;
      }

      tableLen = getData.data.count;
      tableData = getData.data.txs;
      $("#block-table").createTable({ tableData: tableData, tableLength: tableLen, tableId: "#block-table", pageId: "#table-paginate", tableInfoId: "#table-info", sizePage: sizePageParam, currentPage: 1, channelId: channelID, choiceContent: choiceTxt, previewUrl: url })

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      $("#block-load .loading").empty()
      $("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')
    });

}

// get the data of creator of tx
function creatorOftxData(url) {
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
      $("#transaction-creater option:not(:first-child)").remove()

      if (getData.code == 200) {
        arrayData = getData.data
        getDataLen = arrayData.length;

        if (!!getDataLen) {
          for (var i = 0; i < getDataLen; i++) {
            strHtml += '<option value="' + arrayData[i] + '">' + arrayData[i] + '</option>'
          }
          $("#transaction-creater").append(strHtml)
        }
      } else {
        alert(getData.msg)
      }

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      alert("Request failed: " + textStatus);
    });

}
// get the data of tx hash
function txHashTableData(url) {
  loading($(".block-table-contain"))
  $.ajax({
    url: url,
  })
    .done(function (data) {
      $("#block-load").remove();
      $("#block-table tbody tr").remove()
      var getData;

      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }
      if (getData.code !== 200) {
        $("#block-table").append('<tr><td class="text-center" colspan="8">' + getData.msg + '</td></tr>')
        return false;
      }

      getData = getData.data;
      if (!!getData) {
        $("#block-table").append('<tr class="odd"><td>' + getData.block_num + '</td><td>' + getData.tx_num + '</td><td class="block-name"><div class="data-hash" data-toggle="modal" data-target="#myModal2" title="' + getData["tx_hash"] + '">' + getData["tx_hash"] + '</div></td><td><div class="block-hash">' + getData.chaincode_name + '</div></td><td><div class="previous-hash">' + getData.tx_type + '</div></td><td>' + getData.create_msp_id + '</td><td>' + getData.endorser_msp_id.replace(/,/g, '/') + '</td><td>' + getData.created_at.slice(0, 19).replace(/T/g, " ") + '</td></tr>');
      } else {
        $("#block-table").append('<tr><td class="text-center" colspan="8">find no transaction...</td></tr>')
      }

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      $("#block-load .loading").empty()
      $("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')
    });

}
// Tx Hash details
function txHashDetails(url) {
  $.ajax({
    url: url,
  })
    .done(function (data) {
      var getData;

      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }
      getData = getData.data
      var $form = $("#block-form2 .block-detail span")

      $form.eq(0).text(getData.block_num)
      $form.eq(1).text(getData.tx_num)
      $form.eq(2).text(getData.tx_hash)
      $form.eq(3).text(getData.chaincode_name)
      $form.eq(4).text(getData.tx_type)
      $form.eq(5).text(getData.chaincode_function)
      $form.eq(6).text(getData.chaincode_function_parameters)

      $form.eq(7).text(getData.create_msp_id)
      $form.eq(8).text(getData.endorser_msp_id.replace(/,/g, '/'))
      $form.eq(9).text(getData.created_at.slice(0, 19).replace(/T/g, " "))

      var options = {
        collapsed: false,
        withQuotes: true
      };

      $('#read-set').jsonViewer(JSON.parse(getData.read_set), options);

      var writeArr = JSON.parse(getData.write_set)
      if (!writeArr) {
        var writeArrLen = writeArr.length,
          writeSetArr,
          writeSetArrLen;

        for (var i = 0; i < writeArrLen; i++) {
          writeSetArr = writeArr[i]["set"];
          writeSetArrLen = writeSetArr ? writeSetArr.length : 0;
          for (var j = 0; j < writeSetArrLen; j++) {
            writeSetArr[j]["value"] = JSON.parse(writeSetArr[j]["value"])
          }
        }

      }

      $('#write-set').jsonViewer(writeArr, options);

      $("#endorser-signature").jsonViewer(JSON.parse(getData.endorser_signature), options);
    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      alert("Request failed: " + textStatus);
    });
}
