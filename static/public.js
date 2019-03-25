/**
 * Created by aboc on 17-2-24.
 */
var renderData = {
  content:[],
  start:0,
    from:'',
    to:'',
}
const RENDER_NUM = 20;
function resetRenderData() {
    renderData.content = [];
    renderData.start = 0;
    renderData.from ='';
    renderData.to = '';
}
$(function() {
  $.getJSON("/log/", function(d) {
    var str = [];
    for (var i in d) {
      str.push('<li><a href="#" data-commit="' + d[i].Commit + '" title="' + d[i].Commit + '">' + d[i].Author + '<em>' + d[i].Date + '</em></a></li>');
    }
    if (!str.length) {
      str.push('<li><a>没有记录！</a></li>')
    }
    $(".nav-sidebar").html(str.join(""));
  });
  $(".nav-sidebar").delegate("li a", "click", function() {
    resetRenderData();
    var li = $(this).closest("li");
    var next_li = $(li).next();
    var from = $(this).data("commit");
    if (next_li.length) {
      var to = $("a", $(next_li)).data("commit");
    } else {
      var to = "";
    }
    loadView(from,to,0)
    return false;
  });
});

function loadView(from,to,start) {
    layer.load(2)
    console.info(to);
    $.getJSON("/view/", {
        from: from.substr(0, 7),
        to: to.substr(0, 7),
        start:start
    }, function(d) {
        layer.closeAll();
        $(".table-responsive").html('');
        //渲染右侧
        renderData.from = from;
        renderData.to = to;
        renderData.content = d;
        renderData.start = 0;
        startRenderData()
    });
}

function startRenderData() {
    var str = [];
    for(var i in renderData.content) {
      var d = renderData.content[i];
        var log = '<div class="one-file">'
            +'<div class="title">'+d.Filename+'</div>'
            +'<div class="file-log">';
        var num = 1;
        for (var k in d.Lines) {
            var line = d.Lines[k];
            log += '<dl class="'+get_class(line)+'">'
                +'<dt>'+num+'</dt>'
                +'<dd>'+line.replace("<","&lt;").replace(">","&gt;")+'</dd>'
                +'</dl>';
            num ++;
        }
        log +'</dl></div></div>';
        str.push(log);
    }
    var html = str.join("");
    console.info(html);
    $(".table-responsive").html(str.join(""));
}

function get_class(str){
    var f = str.substr(0,1);
  if(f == "+" && str.substr(0,3)!= "+++"){
    return 'jia';
  }
  if(f == "-" && str.substr(0,3)!= "---"){
    return "jian";
  }
  return "";
}