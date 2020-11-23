; (function ($, window, document, undefined) {
	sessionStorage.clear()
	var createTable = function (element, options) {
		this.$element = element;
		this.options = $.extend(true, {
			tableData: [],  // json table data
			tableLength: '',  // total data length
			choiceContent: '',
			pageId: '',
			sizePage: 10, // 10 data per page
			totalPageId: '', //table page id element
			tableId: '',	//table id element
			currentPage: 1, // current page
			channelId: '',
			previewUrl: '',
			tableInfoId: '' // page information id
		}, options);
	};
	createTable.prototype = {

		init: function () {
			this.tableData = this.options.tableData,
				this.tableDataLen = this.options.tableLength,
				this.tableCurrentDataLen = (this.tableData == null) ? 0 : this.tableData.length,
				this.channelId = this.options.channelId,
				this.tablePageLen = '',  // the length of table data
				this.tableStrHtml = '',
				this.choiceContent = this.options.choiceContent,
				this.currentDataStartLen = '',//  the start page of current range
				this.currentDataEndLen = '',
				this.pageStrHtml = '',
				this.prevPage = '',
				this.nextPage = '',
				this.totalPage = '',
				this.currentPage = this.options.currentPage,     //The default current number of pages is page 1
				this.postUrl = this.options.previewUrl,
				this.currPage = '',       //click to get the current number of pages
				this.sizePage = this.options.sizePage,   // 10 data of per page
				this.pageId = this.options.pageId,
				this.totalPageId = this.options.totalPageId,

				this.tableId = this.options.tableId,

				this.initTable(this.currentPage);

			this.clickPage();
		},
		initTable: function (currPageParam) {
			this.totalPage = Math.ceil(this.tableDataLen / this.sizePage);  // caculate the page of data
			this.tablePageLen = this.sizePage;

			$("#block-table tbody tr").remove();

			if (this.tableDataLen !== 0) {				// check if the length of table data is greater than 0
				var i = 0;

				this.currentDataStartLen = (this.currentPage - 1) * this.sizePage + 1;
				this.currentDataEndLen = this.currentDataStartLen + this.tableData.length - 1;

				this.tableStrHtml = ""
				if (this.choiceContent == "blockTime") {
					for (i; i < this.tableCurrentDataLen; i++) {
						this.tableStrHtml += '<tr class="odd"><td>' + this.tableData[i]["block_num"] + '</td><td  class="block-name" data-arrNumber="' + i + '"><div data-toggle="modal" data-target="#myModal" class="data-hash" title="' + this.tableData[i]["data_hash"] + '">' + this.tableData[i]["data_hash"] + '</div></td><td><div class="block-hash" title="' + this.tableData[i]["block_hash"] + '">' + this.tableData[i]["block_hash"] + '</div></td><td><div class="previous-hash" title="' + this.tableData[i]["previous_hash"] + '">' + this.tableData[i]["previous_hash"] + '</div></td><td class="tx-count"><a href="/transaction.html?booknumber=' + this.tableData[i]["block_num"] + '">' + this.tableData[i]["tx_count"] + '</a></td><td>' + this.tableData[i]["invalid_tx_count"] + '</td><td>' + this.tableData[i]["created_at"].slice(0, 19).replace(/T/g, " ") + '</td></tr>';
					}
				}
				if (this.choiceContent == "transactionTime") {
					var transactionMsp;
					for (i; i < this.tableCurrentDataLen; i++) {
						transactionMsp = (this.tableData[i]["endorser_msp_id"] == null) ? "" : this.tableData[i]["endorser_msp_id"].replace(/,/g, "/")

						this.tableStrHtml += '<tr class="odd"><td>' + this.tableData[i]["block_num"] + '</td><td>' + this.tableData[i]["tx_num"] + '</td><td  class="block-name" data-arrNumber="' + i + '"><div data-toggle="modal" data-target="#myModal2" class="data-hash" title="' + this.tableData[i]["tx_hash"] + '">' + this.tableData[i]["tx_hash"] + '</div></td><td><div class="block-hash">' + this.tableData[i]["chaincode_name"] + '</div></td><td><div class="previous-hash">' + this.tableData[i]["tx_type"] + '</div></td><td>' + this.tableData[i]["create_msp_id"] + '</a></td><td><div class="endorser-msp-id">' + transactionMsp + '</div></td><td>' + this.tableData[i]["created_at"].slice(0, 19).replace(/T/g, " ") + '</td></tr>';
					}
				}
				if (this.choiceContent == "transactionBlockNumber") {
					var transactionMsp;
					for (i; i < this.tableCurrentDataLen; i++) {
						transactionMsp = (this.tableData[i]["endorser_msp_id"] == null) ? "" : this.tableData[i]["endorser_msp_id"].replace(/,/g, "/")

						this.tableStrHtml += '<tr class="odd"><td>' + this.tableData[i]["block_num"] + '</td><td>' + this.tableData[i]["tx_num"] + '</td><td  class="block-name" data-arrNumber="' + i + '"><div data-toggle="modal" data-target="#myModal2" class="data-hash" title="' + this.tableData[i]["tx_hash"] + '">' + this.tableData[i]["tx_hash"] + '</div></td><td><div class="block-hash">' + this.tableData[i]["chaincode_name"] + '</div></td><td><div class="previous-hash">' + this.tableData[i]["tx_type"] + '</div></td><td>' + this.tableData[i]["create_msp_id"] + '</a></td><td><div class="endorser-msp-id">' + transactionMsp + '</div></td><td>' + this.tableData[i]["created_at"].slice(0, 19).replace(/T/g, " ") + '</td></tr>';
					}
				}
				if (this.choiceContent == "historyTime") {
					for (i; i < this.tableCurrentDataLen; i++) {
						this.tableStrHtml += '<tr class="odd"><td>' + this.tableData[i]["block_num"] + '</td><td>' + this.tableData[i]["tx_num"] + '</td><td><div class="history-hash" title="' + this.tableData[i]["tx_hash"] + '">' + this.tableData[i]["tx_hash"] + '</div></td><td>' + this.tableData[i]["chaincode_name"] + '</td><td><div class="tx-type">' + this.tableData[i]["tx_type"] + '</div></td><td><div class="create-msp-id">' + this.tableData[i]["create_msp_id"] + '</div></td><td><div class="endorser-msp-id">' + this.tableData[i]["endorser_msp_id"].replace(/,/g, "/") + '</div></td><td class="history-name"><div class="history-key" title="' + this.tableData[i]["key"] + '" data-toggle="modal" data-target="#myModal">' + this.tableData[i]["key"] + '</div><a class="fa fa-list-ul" href="/history.html?historyParams=' + this.tableData[i]["key"] + '" target="_blank" aria-hidden="true" title="see more details"></a></td><td class="schema-position"><div class="history-content">' + this.tableData[i]["content"] + '</div></td><td>' + this.tableData[i]["is_delete"] + '</td><td>' + this.tableData[i]["created_at"].slice(0, 19).replace(/T/g, " ") + '</td></tr>';
					}
				}
				if (this.choiceContent == "historyKey") {
					for (i; i < this.tableCurrentDataLen; i++) {
						this.tableStrHtml += '<tr class="odd"><td>' + this.tableData[i]["block_num"] + '</td><td>' + this.tableData[i]["tx_num"] + '</td><td ><div class="history-hash" title="' + this.tableData[i]["tx_hash"] + '">' + this.tableData[i]["tx_hash"] + '</div></td><td>' + this.tableData[i]["chaincode_name"] + '</td><td><div class="tx-type">' + this.tableData[i]["tx_type"] + '</div></td><td><div class="create-msp-id">' + this.tableData[i]["create_msp_id"] + '</div></td><td><div class="endorser-msp-id">' + this.tableData[i]["endorser_msp_id"].replace(/,/g, "/") + '</div></td><td class="history-name"><div class="history-key" title="' + this.tableData[i]["key"] + '" data-toggle="modal" data-target="#myModal">' + this.tableData[i]["key"] + '</div></td><td class="schema-position"><div class="history-content">' + this.tableData[i]["content"] + '</div></td><td>' + this.tableData[i]["is_delete"] + '</td><td>' + this.tableData[i]["created_at"].slice(0, 19).replace(/T/g, " ") + '</td></tr>';
					}
				}
				if (this.choiceContent == "schema") {

					for (i; i < this.tableCurrentDataLen; i++) {
						this.tableStrHtml += '<tr class="odd">'
						var tableSchemaHtml = ""
						for (var key in this.tableData[i]) {
							if (key == "channel_id" || key == "chaincode_name" || key == "is_delete") {
								tableSchemaHtml += '<td class="display-none"><div>' + this.tableData[i][key] + '</div></td>'
							} else if (key == "doc_key") {
								if (this.isJSON(this.tableData[i][key])) {
									if (typeof (this.tableData[i][key]) == "object") {
										tableSchemaHtml += '<td class="schema-dockey-td"><div class="schema-div"><div class="schema-dockey schema-dockey-object" data-toggle="modal" data-target="#myModal" title="' + JSON.stringify(this.tableData[i][key]) + '">' + JSON.stringify(this.tableData[i][key]) + '</div><a class="fa fa-list-ul" href="/history.html?historyParams=' + JSON.stringify(this.tableData[i][key]) + '" target="_blank" aria-hidden="true" title="see more details"></a></div></td>'
									} else {
										tableSchemaHtml += '<td class="schema-dockey-td"><div class="schema-div"><div class="schema-dockey schema-dockey-object" data-toggle="modal" data-target="#myModal" title="' + JSON.stringify(this.tableData[i][key]) + '">' + JSON.stringify(this.tableData[i][key]) + '</div><a class="fa fa-list-ul" href="/history.html?historyParams=' + this.tableData[i][key] + '" target="_blank" aria-hidden="true" title="see more details"></a></div></td>'
									}
								} else {
									tableSchemaHtml += '<td class="schema-dockey-td"><div class="schema-div"><div class="schema-dockey" data-toggle="modal" data-target="#myModal" title="' + this.tableData[i][key] + '">' + this.tableData[i][key] + '</div><a class="fa fa-list-ul" href="/history.html?historyParams=' + this.tableData[i][key] + '" target="_blank" aria-hidden="true" title="see more details"></a></div></td>'
								}
							} else {
								if (this.isJSON(this.tableData[i][key])) {
									if (typeof (this.tableData[i][key]) == "object") {
										tableSchemaHtml += '<td class="schema-position"><div class="schema-dockey schema-dockey-object schema-copy-object" title="点击可以复制全部内容">' + JSON.stringify(this.tableData[i][key]) + '</div></td>'
									} else {
										tableSchemaHtml += '<td class="schema-position"><div class="schema-dockey schema-dockey-object schema-copy-object" title="点击可以复制全部内容">' + this.tableData[i][key] + '</div></td>'
									}
								} else {
									if (this.tableData[i][key] !== null) {
										if (this.tableData[i][key].toString().length > 12) {
											tableSchemaHtml += '<td><div class="schema-dockey" title="' + this.tableData[i][key] + '">' + this.tableData[i][key] + '</div></td>'
										} else {
											tableSchemaHtml += '<td><div class="schema-dockey">' + this.tableData[i][key] + '</div></td>'
										}
									} else {
										tableSchemaHtml += '<td><div class="schema-dockey">' + this.tableData[i][key] + '</div></td>'
									}

								}
							}

						}


						this.tableStrHtml += tableSchemaHtml

						this.tableStrHtml += "</tr>"

					}

				}

				this.$element.append(this.tableStrHtml);

			} else {
				if (this.choiceContent == "transactionTime" || this.choiceContent == "transactionBlockNumber") {
					$(this.tableId + " tr:not(:first-child)").remove();
					this.$element.append('<tr><td class="text-center" colspan="8">No data is available at this time...</td></tr>')
				} else if (this.choiceContent == "historyTime" || this.choiceContent == "historyKey") {
					$(this.tableId + " tr:not(:first-child)").remove();
					this.$element.append('<tr><td class="text-center" colspan="11">No data is available at this time...</td></tr>')
				} else if (this.choiceContent == "blockTime") {
					$(this.tableId + " tr:not(:first-child)").remove();
					this.$element.append('<tr><td class="text-center" colspan="7">No data is available at this time...</td></tr>')
				}

			}

			this.initPage(currPageParam);
			this.clickPrevPage()
			this.clickNextPage()
		},
		initPage: function (currPageParam) {

			this.pageStrHtml = ""
			$(this.pageId + ">li:not('#example1_next,#example1_previous')").remove();     // page id


			if (this.totalPage <= 1) {
				$("#example1_next").addClass("disabled");
				$(this.pageId + " #example1_previous").after('<li class="paginate_button active"><a>1</a>');
				$("#number-start").text(this.currentDataStartLen)
				$("#number-end").text(this.currentDataEndLen)
				$("#number-total").text(this.tableDataLen)


			} else if (this.totalPage > 1 && this.totalPage < 8) {

				if (currPageParam == 1) {

					if (!$("#example1_previous").hasClass("disabled")) {
						$("#example1_previous").addClass("disabled")
						$("#example1_next").removeClass("disabled")

					}
				} else if (currPageParam == this.totalPage) {
					if (!$("#example1_next").hasClass("disabled")) {
						$("#example1_next").addClass("disabled")
						$("#example1_previous").removeClass("disabled")
					}

				} else {

					$("#example1_previous,#example1_next").removeClass("disabled")

				}

				for (var i = 0; i < this.totalPage; i++) {

					if (currPageParam !== 1) {
						if (i == (currPageParam - 1)) {
							this.pageStrHtml += '<li class="paginate_button active"><a>' + (i + 1) + '</a>'
						} else {
							this.pageStrHtml += '<li class="paginate_button"><a>' + (i + 1) + '</a>'
						}
					} else {
						$("#example1_previous").addClass("disabled");
						if (i == 0) {
							this.pageStrHtml += '<li class="paginate_button active"><a>' + (i + 1) + '</a>'
						} else {
							this.pageStrHtml += '<li class="paginate_button"><a>' + (i + 1) + '</a>'
						}
					}

				}
				$(this.pageId + " #example1_previous").after(this.pageStrHtml);

				$("#number-start").text(this.currentDataStartLen)
				$("#number-end").text(this.currentDataEndLen)
				$("#number-total").text(this.tableDataLen)

			} else if (this.totalPage > 7) {

				if (currPageParam == 1) {
					if (!$("#example1_previous").hasClass("disabled")) {
						$("#example1_previous").addClass("disabled")
						$("#example1_next").removeClass("disabled")

					}
				} else if (currPageParam == this.totalPage) {
					if (!$("#example1_next").hasClass("disabled")) {
						$("#example1_next").addClass("disabled")
						$("#example1_previous").removeClass("disabled")
					}

				} else {
					$("#example1_previous,#example1_next").removeClass("disabled")

				}
				if (currPageParam < 5) {
					for (var i = 0; i < 7; i++) {
						if (currPageParam == i) {

							this.pageStrHtml += '<li class="paginate_button active"><a>' + currPageParam + '</a>'
						}
						if (i == 5) {
							this.pageStrHtml += '<li class="paginate_button disabled" id="example1_ellipsis"><a>…</a></li>'
						}
						if (i == 6) {
							this.pageStrHtml += '<li class="paginate_button"><a>' + this.totalPage + '</a>'
						}

						if (i !== 5 && i !== 6 && i != (currPageParam - 1)) {


							this.pageStrHtml += '<li class="paginate_button"><a>' + (i + 1) + '</a>'
						}

					}
				}

				if (currPageParam > 4 && currPageParam < (this.totalPage - 3)) {


					for (var i = (currPageParam - 3); i < (currPageParam + 4); i++) {

						if (currPageParam == i) {

							this.pageStrHtml += '<li class="paginate_button active"><a>' + currPageParam + '</a>'
						}
						if (i == (currPageParam - 2)) {

							this.pageStrHtml += '<li class="paginate_button disabled" id="example1_ellipsis"><a>…</a></li>'
						}
						if (i == (currPageParam - 3)) {

							this.pageStrHtml += '<li class="paginate_button"><a>1</a>'
						}
						if (i == (currPageParam + 2)) {

							this.pageStrHtml += '<li class="paginate_button disabled" id="example1_ellipsis"><a>…</a></li>'
						}
						if (i == (currPageParam + 3)) {

							this.pageStrHtml += '<li class="paginate_button"><a>' + this.totalPage + '</a>'
						}
						if (i !== (currPageParam - 2) && i !== (currPageParam - 3) && i !== (currPageParam + 2) && i !== (currPageParam + 3) && i !== currPageParam) {

							this.pageStrHtml += '<li class="paginate_button"><a>' + i + '</a>'
						}

					}
				}
				if (currPageParam > (this.totalPage - 4)) {

					for (var i = (this.totalPage - 6); i < (this.totalPage + 1); i++) {

						if (currPageParam == i) {

							this.pageStrHtml += '<li class="paginate_button active"><a>' + currPageParam + '</a>'
						}
						if (i == (this.totalPage - 5)) {
							this.pageStrHtml += '<li class="paginate_button disabled" id="example1_ellipsis"><a>…</a></li>'
						}
						if (i == (this.totalPage - 6)) {
							this.pageStrHtml += '<li class="paginate_button"><a>1</a>'
						}

						if (i !== (this.totalPage - 5) && i !== (this.totalPage - 6) && i !== currPageParam) {

							this.pageStrHtml += '<li class="paginate_button"><a>' + i + '</a>'
						}

					}
				}



				$(this.pageId + " #example1_previous").after(this.pageStrHtml);

				$("#number-start").text(this.currentDataStartLen)
				$("#number-end").text(this.currentDataEndLen)
				$("#number-total").text(this.tableDataLen)

			}


			// this.clickPrevPage();
			// this.clickNextPage();
			// this.clickSkipPage();

		},
		clickPage: function (prevNextParam) {

			var self = this;

			var bookDataLen = self.bookDataLen;
			$(this.pageId).off("click").on("click", "li:not('#example1_previous,#example1_next')", function () {
				var $thisPageIndex = $(this).children().text();
				if (self.totalPage == 1) return false;

				$(this).addClass("active").siblings().removeClass("active");

				self.pageTableContent($thisPageIndex)

				// self.pageTableContent();
				// this.initTable()
				return false;
			})
		},
		clickPrevPage: function () {

			var self = this;

			$("#example1_previous").off("click").on("click", function () {
				var $thisPagePrevIndex = Number($(self.pageId + " .active").children().text())
				if ($thisPagePrevIndex == 1) return false;
				$thisPagePrevIndex = $thisPagePrevIndex - 1;

				self.pageTableContent($thisPagePrevIndex)
				return false;
			})
		},
		clickNextPage: function () {
			var self = this;
			$("#example1_next").off("click").on("click", function () {
				var $thisPageNextIndex = Number($(self.pageId + " .active").children().text());

				if ($thisPageNextIndex == self.totalPage) return false;

				$thisPageNextIndex = $thisPageNextIndex + 1
				self.pageTableContent($thisPageNextIndex)
				return false;

			})
		},
		pageTableContent: function ($thisPageIndex) {
			var self = this;
			var postParams = {}
			if (this.choiceContent == "blockTime") {
				postParams.channel_id = self.channelId
				postParams.sort_by = $("#block-table th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC";
				postParams.start_time = $("#start-time").val().trim()
				postParams.end_time = $("#end-time").val().trim()
				postParams.page_num = Number($thisPageIndex)
				postParams.page_count = Number($("#select-paging").val())

			} else if (this.choiceContent == "transactionTime") {

				postParams.channel_id = self.channelId
				postParams.sort_by = $("#block-table tr th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC"
				postParams.start_time = $("#start-time").val().trim()
				postParams.end_time = $("#end-time").val().trim()
				postParams.page_num = Number($thisPageIndex)
				postParams.page_count = Number($("#select-paging").val())
				postParams.create_msp = $("#transaction-creater").val() == "null" ? "" : $("#transaction-creater").val()


			} else if (this.choiceContent == "transactionBlockNumber") {

				postParams.channel_id = self.channelId
				postParams.page_num = Number($thisPageIndex)
				postParams.page_count = Number($("#select-paging").val())
				postParams.block_num = Number($("#block-number-input").val().trim())
				postParams.sort_by = $("#block-table tr th:eq(1)").hasClass("sorting_desc") ? "DESC" : "ASC"

			} else if (this.choiceContent == "historyTime") {
				postParams.channel_id = self.channelId
				postParams.sort_by = $("#block-table thead tr th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC"
				postParams.start_time = $("#start-time").val().trim()
				postParams.end_time = $("#end-time").val().trim()
				postParams.page_num = Number($thisPageIndex)
				postParams.page_count = Number($("#select-paging").val())

			} else if (this.choiceContent == "historyKey") {
				postParams.channel_id = self.channelId
				postParams.sort_by = $("#block-table thead tr th:eq(0)").hasClass("sorting_desc") ? "DESC" : "ASC"
				postParams.doc_key = $("#block-number-input").val().trim()
				postParams.create_msp = ($("#transaction-creater").val() == "null") ? "" : $("#transaction-creater").val()
				postParams.chaincode_name = $("#chaincode-name-dockey").val() == "null" ? "" : $("#chaincode-name-dockey").val()
				postParams.page_num = Number($thisPageIndex)
				postParams.page_count = Number($("#select-paging").val())

			} else if (this.choiceContent == "schema") {

				var schemaLeftArr = [], schemaRightArr = [], regTest = /^[-]?[0-9]+\.?[0-9]?$/,
					$selectLeftLen = $("#schema-left-input .select-left-item").length,
					$selectRightLen = $("#schema-right-input .select-right-item").length;

				$("#schema-left-input .select-left-item").each(function (index) {
					if (index == ($selectLeftLen - 1)) return false;
					var $selectionInput = $(this).children(),
						$thisInput1 = $selectionInput.eq(0).children().val().trim(),
						$thisInput2 = $selectionInput.eq(2).children().val().trim();
					if (!($thisInput1 == "" && $thisInput2 == "")) {
						schemaLeftArr.push("'{" + $thisInput1 + "}' AS " + $thisInput2)
					}

				})

				$("#schema-right-input .select-right-item").each(function (index) {
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

							if (!reg.test($thisInput2)) {
								$thisInput2 = "'" + $thisInput2 + "'"
							} else {
								$thisInput2 = Number($thisInput2)
							}

						}
						$thisInput1 = "'{" + $thisInput1 + "}'"
						schemaRightArr.push($thisInput1 + " " + $selectVal + " " + $thisInput2)
					}


				})

				postParams.channel_id = self.channelId
				postParams.schema_id = Number($(".schema-active").attr("data-schemaId"))
				postParams.selects = schemaLeftArr
				postParams.conditions = schemaRightArr
				postParams.page_num = $thisPageIndex
				postParams.page_count = Number($("#select-paging").val())

			}

			//loading html
			$(".block-table-contain").append('<div class="loading-bg" id="block-load"><div class="loading"><span></span><span></span><span></span><span></span><span></span><span></span><span></span></div></div>')

			$.ajax({
				url: this.postUrl,
				method: "POST",
				data: JSON.stringify(postParams),
				contentType: 'application/json',
				dataType: "json"
			})
				.done(function (data) {
					$("#block-load").remove();
					var getData,
						prototypeName;
					if (typeof (data) == "string") {
						getData = JSON.parse(data);
					} else if (typeof (data) == "object") {
						getData = data;
					}

					switch (self.choiceContent) {
						case "blockTime":
							prototypeName = "blocks";
							break;
						case "transactionTime":
							prototypeName = "txs";
							break;
						case "transactionBlockNumber":
							prototypeName = "txs";
							break;
						case "historyTime":
							prototypeName = "history";
							break;
						case "historyKey":
							prototypeName = "history";
							break;
					}


					self.tableDataLen = getData.data.count
					self.tableData = getData.data[prototypeName]

					self.tableCurrentDataLen = self.tableData.length;

					self.currentPage = Number($thisPageIndex)
					self.totalPage = Math.ceil(self.tableDataLen / self.sizePage)

					self.initTable(self.currentPage);
				})
				.fail(function (jqXHR, textStatus) {
					networkStatus(jqXHR.status)
					$("#block-load .loading").empty()
					$("#block-load .loading").append('<div class="text-center error-color">Error: Check Your Inputs Please， And Refresh This Page!</div>')
				});


		},
		isJSON: function (str) {
			if (typeof str == 'string') {
				try {
					if (str.indexOf('{') > -1) {
						return true;
					} else {
						return false;
					}
				}
				catch (e) {
					console.log(e);
					return false;
				}
			} else if (typeof str == 'object') {

				try {

					if (JSON.stringify(str).indexOf('{') > -1) {
						return true;
					} else {
						return false;
					}
				}
				catch (e) {
					console.log(e);
					return false;
				}
			}
		}


	};
	$.fn.createTable = function (options) {

		var obj = new createTable(this, options);

		obj.init();

		return this;
	};


})(jQuery, window, document);
