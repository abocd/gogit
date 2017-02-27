/**
 * Created by aboc on 17-2-24.
 */

$(function() {
  $.getJSON("/log/", function(d) {
    var str = [];
    for (var i in d) {
      str.push('<li><a href="#" data-commit="' + d[i].Commit + '">' + d[i].Date + '</a></li>');
    }
    if (!str.length) {
      str.push('<li><a>没有记录！</a></li>')
    }
    $(".nav-sidebar").html(str.join(""));
  });
  $(".nav-sidebar").delegate("li a", "click", function() {
    var li = $(this).closest("li");
    var next_li = $(li).next();
    var from = $(this).data("commit");
    if (next_li.length) {
      var to = $("a", $(next_li)).data("commit");
    } else {
      var to = "";
    }
    console.info(to);
    $.getJSON("/view/", {
      from: from.substr(0, 7),
      to: to.substr(0, 7)
    }, function(d) {
        //渲染右侧
      var str = [];
      for(var i in d) {
        var log = '<div class="one-file">'
        +'<div class="title">'+d[i].FileName+'</div>'
        +'<div class="file-log">';
        var num = 1;
        for (var k in d[i].Lines) {
          var line = d[i].Lines[k];
              log += '<dl class="'+get_class(line)+'">'
                     +'<dt>'+num+'</dt>'
                     +'<dd>'+line+'</dd>'
                     +'</dl>';
              num ++;
        }
        log +'</dl></div></div>';
        str.push(log);
      }
      console.info(str)
      $(".table-responsive").html(str.join(""));
    });
    return false;
  });
});


function get_class(str){
    var f = str.substr(0,1);
  if(f == "+"){
    return 'jia';
  }
  if(f == "-"){
    return "jian";
  }
  return "";
}