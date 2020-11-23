/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
$(function () {
    $("input").val("")
    $("#schema-right-input select.form-control").val("=")
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

            // parameters of document schemas
            // chaincodeNameParam.channelId = channelID

            // get data of chaincode Name
            chainCodeNameData(chaincodeNameUrl + channelID, channelID)

            //choose schema item
            $("#json-content").off("click").on("click", ".schema", function () {
                $(this).addClass("schema-active").parent().siblings().children().removeClass("schema-active")
            })

            //the first line of left part can't be deleted
            $("#schema-left-input").on("click", ".schema-close", function () {
                if ($("#schema-left-input .schema-close").length == 1) {
                    return false;
                }
                $(this).parent().parent().parent().remove()
            })

            //the first line of Right part can't be deleted
            $("#schema-right-input").on("click", ".schema-close", function () {
                if ($("#schema-right-input .schema-close").length == 1) {
                    return false;
                }
                $(this).parent().parent().parent().remove()
            })

            // click 'search' button , get some data
            $("#schema-search").on("click", debounce(function () {
                $("#block-load").remove();
                $("#schema-table-section").addClass("display-none")
                var schemaBoolean = true;

                if (!$(".schema").hasClass("schema-active")) { //  check if a schema is choosed
                    alert("select document schemas first")
                    schemaBoolean = false;
                }
                if (!schemaBoolean) return schemaBoolean

                var $leftItmeLen = $("#schema-left-input .select-left-item").length,
                    $rightItmeLen = $("#schema-right-input .select-right-item").length;

                // make sure that each row of data of left part is filled.
                $("#schema-left-input .select-left-item").each(function (index) {
                    var $inputList = $(this).children(),
                        $thisInput1 = $inputList.eq(0).children().val().trim(),
                        $thisInput2 = $inputList.eq(2).children().val().trim();

                    if (index == ($leftItmeLen - 1)) {
                        if (!(($thisInput1 == "" && $thisInput2 == "") || ($thisInput1 !== "" && $thisInput2 !== ""))) {
                            alert("Please fill in " + (index + 1) + " lines of data in SELECTIONS")
                            schemaBoolean = false;
                            return false;
                        }
                    } else {

                        if (!($thisInput1 !== "" && $thisInput2 !== "")) {
                            alert("Please fill in " + (index + 1) + " lines of data in SELECTIONS")
                            schemaBoolean = false;
                            return false;
                        }
                    }

                })
                if (!schemaBoolean) return schemaBoolean

                // make sure that each row of data of right part is filled.
                $("#schema-right-input .select-right-item").each(function (index) {
                    var $inputList = $(this).children(),
                        $thisInput1 = $inputList.eq(0).children().val().trim(),
                        $thisInput2 = $inputList.eq(2).children().val().trim();
                    if (index == ($rightItmeLen - 1)) {

                        if (!(($thisInput1 == "" && $thisInput2 == "") || ($thisInput1 !== "" && $thisInput2 !== ""))) {

                            alert("Please fill in " + (index + 1) + " lines of data in CONDITIONS")
                            schemaBoolean = false;
                            return false;
                        }
                    } else {

                        if (!($thisInput1 !== "" && $thisInput2 !== "")) {
                            alert("Please fill in " + (index + 1) + " lines of data in CONDITIONS")
                            schemaBoolean = false;
                            return false;
                        }
                    }

                })
                if (!schemaBoolean) return schemaBoolean

                // the validation passed
                $("#schema-table-section").hasClass("display-none") ? "" : $("#schema-table-section").addClass("display-none")

                searchResultsTableData(advanceTableUrl, channelID, "schema")

            }))

            $("#select-paging").on("change", function () {
                searchResultsTableData(advanceTableUrl, channelID, "schema")
            })
        })
        .fail(function (jqXHR, textStatus) {
            alert("Request failed: " + textStatus);
        });

    // The number of rows on the form that are automatically added on the left.
    $("#schema-left-input").on("input propertychange", "input", function () {
        var totalVal = '',
            $inputParent = $(this).parent().parent(),
            $input = $inputParent.children(),
            $inputNext = $inputParent.next();

        totalVal = $input.eq(0).children().val().trim() + $input.eq(2).children().val().trim();

        if (totalVal == "") {
            if ($inputNext.attr("class") == "select-left-item") {
                $inputNext.remove()
            }
        } else {
            if ($inputNext.attr("class") !== "select-left-item") {
                $inputParent.after('<div class="select-left-item"><div class="col-sm-5"><input type="text" class="form-control" placeholder="Country,Province,City"></div><div class="col-sm-1 select-left-as"><div class="text-center"> AS </div></div><div class="col-sm-5"><input type="text" class="form-control" placeholder="City"></div><div class="col-sm-1 select-left-close"><div class="text-center"><span class="schema-close">×</span></div></div></div>')
            }
        }

    })
    // The number of rows on the form that are automatically added on the right.
    $("#schema-right-input").on("input propertychange", "input", function () {
        var totalVal = '',
            $inputParent = $(this).parent().parent(),
            $input = $inputParent.children(),
            $inputNext = $inputParent.next();

        totalVal = $input.eq(0).children().val().trim() + $input.eq(2).children().val().trim();

        if (totalVal == "") {
            if ($inputNext.attr("class") == "select-right-item") {
                $inputNext.remove()
            }
        } else {
            if ($inputNext.attr("class") !== "select-right-item") {
                $inputParent.after('<div class="select-right-item"><div class="col-sm-5"><input type="text" class="form-control" placeholder="Country,Province,City"></div><div class="col-sm-1 select-right-select"><select class="form-control" name=""><option value="=">=</option><option value=">=">>=</option><option value=">">></option><option value="<"><</option><option value="<="><=</option><option value="!=">!=</option><option value="LIKE">LIKE</option></select></div><div class="col-sm-5"><input type="text" class="form-control" placeholder="Suzhou"></div><div class="col-sm-1 select-right-close"><div class="text-center"><span class="schema-close">×</span></div></div></div></div>')
            }
        }

    })

    // click the word 'doc_key' in the table, show some details
    $("#block-table").off("click").on("click", ".schema-dockey-td", function () {

        $("#block-form").empty()

        var dataLen = $("#block-table thead th").length,
            strHtml = '';
        $("#block-table thead th").each(function (index) {
            strHtml += '<div class="form-group"><label class="col-sm-3 control-label">' + $(this).text() + '</label><div class="col-sm-9"><div class="block-detail"><span></span></div></div></div>'
        })
        $("#block-form").append(strHtml)

        var options = {
            collapsed: false,
            withQuotes: true
        };
        var $objData = $(this).parent().children(),
            $blockDetails = $("#block-form .block-detail");

        for (var i = 0; i < dataLen; i++) {

            if ($objData.eq(i).children().hasClass("schema-dockey-object")) {
                $blockDetails.eq(i).children().addClass("schema-json")
                $blockDetails.eq(i).children().jsonViewer(JSON.parse($objData.eq(i).children().text()), options);
            } else {
                $blockDetails.eq(i).children().text($objData.eq(i).children().text())
            }
        }


    })

    // click on the content of doc_value to implement the copy function
    copy2($("#schema-dockey"), ".schema-copy-object")

    // clean form
    $("#reset-input").off("click").on("click", function () {
        $("#schema-left-input .select-left-item:not(:first-child),#schema-right-input .select-right-item:not(:first-child)").remove()
        $("#schema-left-input .form-control,#schema-right-input .form-control").val("")
        $("#schema-right-input .select-right-select .form-control").val("=")
    })
})


// get chain code name data
function chainCodeNameData(url, channelID) {

    $.ajax({
        url: url,
        // method: "POST"
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

            arrayData = getData.data
            getDataLen = arrayData.length

            $("#advance-chaincode-name option").remove()
            for (var i = 0; i < getDataLen; i++) {
                strHtml += '<option value="' + arrayData[i] + '">' + arrayData[i] + '</option>';
            }

            $("#advance-chaincode-name").append(strHtml)
            schemaData(advanceSchemaUrl + channelID + "/" + (($("#advance-chaincode-name").val() == "Please select") ? "" : $("#advance-chaincode-name").val()))

            // when chaincodeName changes, get the data back.
            $("#advance-chaincode-name").on("change", function () {

                schemaData(advanceSchemaUrl + channelID + "/" + (($("#advance-chaincode-name").val() == "Please select") ? "" : $("#advance-chaincode-name").val()))
            })
        })
        .fail(function (jqXHR, textStatus) {
            alert("Request failed: " + textStatus);
        });

}

//get Document Schemas data
function schemaData(url) {
    console.log(url)
    loading($(".schema-contain"))
    $.ajax({
        url: url,
        contentType: 'application/json',
        dataType: "json"
        // method: "POST"
    })
        .done(function (data) {
            $("#block-load").remove();
            var getData,
                arrayData,
                getDataLen,
                preHtml = "";
            if (typeof (data) == "string") {
                getData = JSON.parse(data);
            } else if (typeof (data) == "object") {
                getData = data;
            }

            $("#json-content").empty()

            if (getData.code !== 200) { // no data

                $("#json-content").append('<div class="text-center error-color schema-no-channel">find no schemas for this channel</div>')

                return false;
            }
            arrayData = getData.data;
            getDataLen = arrayData.length;

            for (var i = 0; i < getDataLen; i++) {
                preHtml += '  <div class="col-sm-6 col-lg-3 col-md-4 "><div data-schemaId="' + arrayData[i]["id"] + '" class="schema"><input type="text" class="display-none data-input"><div class="schema-title">Schema' + arrayData[i]["id"] + '</div> <pre class="advance-json"></pre></div></div>'
            }
            $("#json-content").append(preHtml)

            for (var i = 0; i < getDataLen; i++) {
                $("#json-content .data-input:eq(" + i + ")").val(JSON.stringify(arrayData[i]["schema_json"]))
            }

            var options = {
                collapsed: false,
                withQuotes: true
            };
            for (var i = 0; i < getDataLen; i++) {
                $("#json-content .advance-json:eq(" + i + ")").jsonViewer(JSON.parse(arrayData[i]["schema_json"]), options);
            }

            if (getDataLen <= 4) {  // the number of schemas less than 4,
                $("#next-slide .next-icon-right,#prev-slide .prev-icon-left").addClass("display-none")
                return false;
            }
            $("#next-slide .next-icon-right").removeClass("display-none")

            //the slides of schemas slides are loaded
            // click next button, move three to the right
            $("#next-slide .next-icon-right").on("click", function () {
                var nextIndex = Number($("#next-slide").attr("data-nextIndex")),
                    prevIndex = Number($("#prev-slide").attr("data-prevIndex")),
                    widthTransation = $(".json-content .col-md-4").outerWidth(true)
                $("#prev-slide .prev-icon-left").removeClass("display-none")

                if (nextIndex == (getDataLen + 1)) return false;

                if ((nextIndex + 2) < getDataLen) {  // When the number of pages on the next page is greater than 3.
                    $(this).removeClass("display-none")
                    $(".json-content").css({
                        "transform": "translate3d(-" + widthTransation * (nextIndex - 2) + "px, 0px, 0px)",
                        "transition": "all .5s"
                    })
                    $("#prev-slide").attr("data-prevIndex", prevIndex + 3)
                    $("#next-slide").attr("data-nextIndex", nextIndex + 3)
                } else { // When the number of pages on the next page is less than 3.

                    $(".json-content").css({
                        "transform": "translate3d(-" + widthTransation * (getDataLen - 4) + "px, 0px, 0px)",
                        "transition": "all .5s"
                    })
                    $("#prev-slide").attr("data-prevIndex", getDataLen - 5)
                    $("#next-slide").attr("data-nextIndex", getDataLen + 1)

                    $(this).addClass("display-none")
                }

            })
            // click prev button, move three to the left
            $("#prev-slide .prev-icon-left").on("click", function () {
                var nextIndex = Number($("#next-slide").attr("data-nextIndex")),
                    prevIndex = Number($("#prev-slide").attr("data-prevIndex")),
                    widthTransation = $(".json-content .col-md-4").outerWidth(true)
                $("#next-slide .next-icon-right").removeClass("display-none")


                if (prevIndex == -1) return false;

                if (prevIndex > 2) {
                    $(".json-content").css({
                        "transform": "translate3d(-" + widthTransation * (prevIndex - 2) + "px, 0px, 0px)",
                        "transition": "all .5s"
                    })
                    $("#prev-slide").attr("data-prevIndex", prevIndex - 3)
                    $("#next-slide").attr("data-nextIndex", nextIndex - 3)
                    $(this).removeClass("display-none")
                } else {
                    $(".json-content").css({
                        "transform": "translate3d(-" + widthTransation * 0 + "px, 0px, 0px)",
                        "transition": "all .5s"
                    })
                    $("#prev-slide").attr("data-prevIndex", -1)
                    $("#next-slide").attr("data-nextIndex", 5)
                    $(this).addClass("display-none")
                }


            })
        })
        .fail(function (jqXHR, textStatus) {
            $("#block-load").remove();
            alert("Request failed: " + textStatus);
        });

}

// query by selections or conditions
//click search button, show details
function searchResultsTableData(url, channelID, choiceContentParams) {

    $("#no-result").remove();
    $("#block-table thead,#block-table tbody").empty()

    var schemaLeftArr = [], schemaRightArr = [], regTest = /^[-]?[0-9]+\.?[0-9]?$/,
        $selectLeftElement = $("#schema-left-input .select-left-item"),
        $selectRightElement = $("#schema-right-input .select-right-item"),
        $selectLeftLen = $selectLeftElement.length,
        $selectRightLen = $selectRightElement.length;
    //selections
    $selectLeftElement.each(function (index) {
        if (index == ($selectLeftLen - 1)) return false;

        var $selectionInput = $(this).children(),
            $thisInput1 = $selectionInput.eq(0).children().val().trim(),
            $thisInput2 = $selectionInput.eq(2).children().val().trim();

        if (!($thisInput1 == "" && $thisInput2 == "")) {
            schemaLeftArr.push("'{" + $thisInput1 + "}' AS " + $thisInput2)
        }

    })
    //conditions
    $selectRightElement.each(function (index) {
        if (index == ($selectRightLen - 1)) return false;

        var $selectionInput = $(this).children(),
            $thisInput1 = $selectionInput.eq(0).children().val().trim(),
            $thisInput2 = $selectionInput.eq(2).children().val().trim(),
            $selectVal = $selectionInput.eq(1).children().val().trim(),
            reg = /^(-?\d+)(\.)?(\d+)?$/g;

        if (!($thisInput1 == "" && $thisInput2 == "")) {

            if ($selectVal == "LIKE") {
                $thisInput2 = "'%" + $thisInput2 + "%'"
            } else {
                if (!reg.test($thisInput2)) { // check if the valuse show 'string'
                    $thisInput2 = "'" + $thisInput2 + "'"
                } else {
                    $thisInput2 = Number($thisInput2)
                }
            }

            $thisInput1 = "'{" + $thisInput1 + "}'"
            schemaRightArr.push($thisInput1 + " " + $selectVal + " " + $thisInput2)
        }

    })

    var schemaParams = {};
    schemaParams.channel_id = channelID
    schemaParams.schema_id = Number($(".schema-active").attr("data-schemaId"))
    schemaParams.selects = schemaLeftArr
    schemaParams.conditions = schemaRightArr
    schemaParams.page_num = 1
    schemaParams.page_count = Number($("#select-paging").val())

    $("#schema-table-section").removeClass("display-none")

    loading($(".block-table-contain"))

    //get table data
    $.ajax({
        url: url,
        timeout: 3000,
        contentType: 'application/json;charset=UTF-8',
        dataType: "json",
        data: JSON.stringify(schemaParams),
        method: "POST"
    })
        .done(function (data) {
            $("#block-load").remove();

            var getData,
                getDataLen,
                getResultData,
                getTableLen,
                theadHtml = '',
                sizePageParam;
            if (typeof (data) == "string") {
                getData = JSON.parse(data);
            } else if (typeof (data) == "object") {
                getData = data;
            }

            if (getData.code !== 200) {
                $("#block-table tbody").empty().append('<div class="text-center error-color">' + getData.msg + '</div>')
                return false;
            }

            getTableLen = getData.data.count
            getResultData = getData.data.result

            $("#block-table thead,#block-table tbody").empty()

            if (getTableLen == 0) {
                $("#schema-table-section").addClass("display-none")
                $("#schema-table-section").before('<div id="no-result" class="no-result">find no results</div>')
            } else {
                getDataLen = getResultData.length
                // table header
                theadHtml += '<tr>'
                for (var key in getResultData[0]) {
                    if (key == "doc_key") {
                        theadHtml += '<th class="th-doc-key">' + key + '</th>'
                    } else {
                        if (key == "channel_id" || key == "chaincode_name" || key == "is_delete") {
                            theadHtml += '<th class="display-none">' + key + '</th>'
                        } else {
                            theadHtml += '<th>' + key + '</th>'
                        }

                    }
                }
                theadHtml += '</tr>'

                $("#block-table thead").append(theadHtml)

                sizePageParam = $("#select-paging").val()

                $("#no-result").remove();
                $("#block-table").createTable({ tableData: getResultData, tableLength: getTableLen, tableId: "#block-table", pageId: "#table-paginate", tableInfoId: "#table-info", sizePage: $("#select-paging").val(), currentPage: 1, channelId: channelID, choiceContent: choiceContentParams, previewUrl: url })

                var scrollheight = $("#schema-table-section").offset().top;
                $("body,html").animate({ scrollTop: scrollheight }, 1000);

            }

        })
        .fail(function (jqXHR, textStatus) {
            console.log("Request failed: " + textStatus);
            $("#block-load .loading").empty()
            $("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')
        });
}
