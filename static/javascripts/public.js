/*
Copyright Fujitsu Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

$(function () {

  if (!!$("#start-time").length && !!$("#end-time").length) {

    $('#start-time').datetimepicker({
      format: 'yyyy-mm-dd hh:ii:ss',
      autoclose: true, // auto off
      endDate: new Date(), //Dates after the current time are not optional
      minView: 2, //The most accurate time selection is the date 0-min 1-hour 2-day 3-month.
      weekStart: 0
    }).on('changeDate', function (ev) {

      var starttime = $("#start-time").val().trim();
      $("#end-time").datetimepicker('setStartDate', starttime);
      $("#start-time").datetimepicker('hide');

    });

    $('#end-time').datetimepicker({
      format: 'yyyy-mm-dd hh:ii:ss',
      autoclose: true,
      endDate: new Date(),
      minView: 2,
      weekStart: 0
    }).on('changeDate', function (ev) {

      var endtime = $("#end-time").val().trim();
      $("#start-time").datetimepicker('setEndDate', endtime);
      $("#end-time").datetimepicker('hide')

    });

  }
})


function loading($element) {
  $element.append('<div class="loading-bg" id="block-load"><div class="loading"><span></span><span></span><span></span><span></span><span></span><span></span><span></span></div></div>')
}
// click icon, you can copy this content
function copy($elementParent, $ele) {
  $elementParent.off("click").on("click", $ele, function () {
    var $elem = $(this).prev(),
      copyHtml = '<i class="fa fa-copy" aria-hidden="true"></i>';
    $(".oInput").remove();

    var Url2 = $(this).prev().text(),
      oInput = document.createElement('input');

    oInput.value = Url2;
    document.getElementById("copy-input").appendChild(oInput);

    oInput.select(); //select this element
    document.execCommand("Copy"); // execute the browser copy command
    oInput.className = 'oInput';
    oInput.style.display = 'none';

    $(this).remove()
    $elem.after('<span class="copyed-style"><i class="fa fa-check-circle"></i><i>&nbsp;&nbsp;copyed</i></span>')

    setTimeout(function () {
      $elem.next().remove();
      $elem.after(copyHtml)
    }, 500)

  })
}
// click doc content or doc_value, you can copy this content
function copy2($elementParent, $ele) {

  $elementParent.off("click").on("click", $ele, function () {

    var $elem = $(this);

    $(".oInput").remove();
    var Url2 = $(this).text();
    var oInput = document.createElement('input');
    oInput.value = Url2;

    document.getElementById("hidden-copy-input").appendChild(oInput);
    oInput.select(); // select Object
    document.execCommand("Copy"); // Execute the browser copy command
    oInput.className = 'oInput';
    oInput.style.display = 'none';

    $elem.after('<span class="copyed-style advance-copy"><i class="fa fa-check-circle"></i><i>&nbsp;&nbsp;copyed</i></span>')
    setTimeout(function () {
      $elem.next().remove();
    }, 500)
  })
}

function networkStatus(status) {
  switch (status) {
    case 404:
      window.location.href = "error4.html";
      break;
    case 503:
      window.location.href = "error5.html";
      break;

  }
}

function debounce(fn) {
  let timeout = null;
  return function () {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      fn.apply(this, arguments);
    }, 300);
  };
}
