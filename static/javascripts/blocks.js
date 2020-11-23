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
      //initialize the table
      initTableData(blockTimeUrl, channelID, "blockTime")

      // search by Time Range
      $("#block-time-search").on("click", debounce(function () {
        initTableData(blockTimeUrl, channelID, "blockTime")
      }))

      // search by Blcok Height
      $("#block-number-search").on("click", debounce(function () {
        var reg = /^\d+$/g
        var $input = $("#block-number-input").val().trim();

        if (!reg.test($input)) {
          alert("Please fill in non-negative integers!")
          return false;
        }

        blockHeightTableData(blockNumberUrl + channelID + "/" + $input)

      }))
      // search by Data Hash
      $("#block-hash-search").on("click", debounce(function () {

        var $input = $("#block-hash-input").val().trim();
        if ($input == "") {
          alert("Please fill in non-negative integers!")
          return false;
        }

        blockHeightTableData(blockDataHashUrl + channelID + "/" + $input)

      }))

      //  select the data, show the number of informations on each page.
      $("#select-paging").on("change", function () {
        initTableData(blockTimeUrl, channelID, "blockTime")
      })

      // sort by Block Height
      $("#block-table th:eq(0)").on("click", debounce(function () {
        var self = $(this);
        (self.hasClass("sorting_desc")) ? self.removeClass("sorting_desc").addClass("sorting_asc") : self.removeClass("sorting_asc").addClass("sorting_desc");

        initTableData(blockTimeUrl, channelID, "blockTime")
      }))

      // click the word 'datahash' in the table, show some details
      $("#block-table").off("click").on("click", ".data-hash", function () {

        var $elem = $(this).parent().parent().children()
        $form = $("#block-form .block-detail span")

        for (var i = 0; i < $form.length; i++) {
          $form.eq(i).text($elem.eq(i).text())
        }

      })

      // click on the icon to implement the copy function
      copy($("#block-form"), ".fa-copy")

      // click the tab element, select the corresponding content part.
      $("#tab-click").find("li").on("click", function () {

        var $thisIndex = $(this).index(),
          $sortElement = $("#block-table thead tr th:eq(0)"),
          $pageElement = $("#paging-length,#entries-info,#entries-page");

        $("#block-load").remove()
        $("#table-search").addClass("animated flipInX")
        setTimeout(function () { $("#table-search").removeClass("animated flipInX") }, 1000)  // dynamic effects
        $("#table-search>.form-inline").eq($thisIndex).removeClass("display-none").siblings().addClass("display-none")
        $(this).addClass("blocks-active").siblings().removeClass("blocks-active")

        $("#start-time,#end-time,#block-number-input,#block-hash-input").val("")  // clear the content of all inputs
        $("#block-table tbody tr").remove()

        $sortElement.removeClass("sorting_desc sorting_asc")

        if ($thisIndex == 0) {  // Time part

          $sortElement.addClass("sorting_desc")
          $pageElement.removeClass("display-none")  // show drop-down option and pages
          initTableData(blockTimeUrl, channelID, "blockTime")

        } else if ($thisIndex == 1) {  // Blcok Height part

          $pageElement.addClass("display-none")
          $("#block-table").append('<tr><td class="text-center" colspan="7">Input Block Height to Query Block Please...</td></tr>')

        } else if ($thisIndex == 2) {  // Data Hash part

          $pageElement.addClass("display-none")
          $("#block-table").append('<tr><td class="text-center" colspan="7">Input Data Hash of Block to Query Block Please...</td></tr>')

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
* The table of Time , table codes.
* use AJAX to get the data.
*/

function initTableData(url, channelID, choiceTxt) {
  $("#block-load").remove();
  loading($(".block-table-contain"))

  var postParams = {},
    sizePageParam = Number($("#select-paging").val());

  postParams.channel_id = channelID
  postParams.sort_by = $("#block-table th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC";
  postParams.start_time = $("#start-time").val().trim()
  postParams.end_time = $("#end-time").val().trim()
  postParams.page_num = 1
  postParams.page_count = sizePageParam

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
      tableData = getData.data.blocks;

      $("#block-table").createTable({ tableData: tableData, tableLength: tableLen, tableId: "#block-table", pageId: "#table-paginate", tableInfoId: "#table-info", sizePage: sizePageParam, currentPage: 1, channelId: channelID, choiceContent: choiceTxt, previewUrl: url })

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      $("#block-load .loading").empty()
      $("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')
    });

}

/*
* The table of Block Height and Data Hash, republic table codes
* use AJAX to get the data
* @param {String} url required parameters
* @param {Object} params required parameters
*/

function blockHeightTableData(url) {
  $("#block-load").remove();
  loading($(".block-table-contain"))
  $.ajax({
    url: url,
  })
    .done(function (data) {

      $("#block-load").remove();
      $("#block-table tbody tr").remove()
      var getData,
          arrayData;

      if (typeof (data) == "string") {
        getData = JSON.parse(data);
      } else if (typeof (data) == "object") {
        getData = data;
      }

      if (getData.code !== 200) {
        $("#block-table").append('<tr><td class="text-center" colspan="7">' + getData.msg + '</td></tr>');
        return false;
      }

      arrayData = getData.data;

      if (!!arrayData) {
        $("#block-table").append('<tr class="odd"><td>' + arrayData.block_num + '</td><td class="block-name"><div class="data-hash" data-toggle="modal" data-target="#myModal" title="' + arrayData["data_hash"] + '">' + arrayData["data_hash"] + '</div></td><td><div class="block-hash" title="' + arrayData.block_hash + '">' + arrayData.block_hash + '</div></td><td><div class="previous-hash" title="' + arrayData.previous_hash + '">' + arrayData.previous_hash + '</div></td><td class="tx-count"><a href="/transaction.html?booknumber=' + arrayData.block_num + '">' + arrayData.tx_count + '</a></td><td>' + arrayData.invalid_tx_count + '</td><td>' + arrayData.created_at.slice(0, 19).replace(/T/g, " ") + '</td></tr>');
      } else {
        $("#block-table").append('<tr><td class="text-center" colspan="7">find no blocks...</td></tr>')
      }

    })
    .fail(function (jqXHR, textStatus) {
      networkStatus(jqXHR.status)
      $("#block-load .loading").empty()
      $("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')

    });


}
