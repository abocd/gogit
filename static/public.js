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
    renderData.startEnd = 0;
    renderData.from ='';
    renderData.to = '';
}
var isLoad = false
var canLoad = true
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
    $(".table-responsive").html('');
    loadView(from,to,0)
    return false;
  });

  $(window).scroll(function (d) {
    if(!isLoad && canLoad){
      if( ($(window).scrollTop() + $(window).height()) > $("body").height() ){
        isLoad = true;
        loadView(renderData.from,renderData.to,renderData.startEnd);
      }
    }
  })
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
        //渲染右侧
        renderData.from = from;
        renderData.to = to;
        renderData.content = d.Files;
        renderData.start = d.Start;
        renderData.startEnd = d.StartEnd;
        if(d.Files != null){
          canLoad = true
        } else {
          canLoad = false
        }
        startRenderData()
    });
}

function startRenderData() {
    var str = [];
    for(var i in renderData.content) {
      var d = renderData.content[i];
        if($(".file_"+d.FileMd5).length > 0){

          var num =  $(".file-log dl",$(".file_"+d.FileMd5)).length;
          console.info("插入到 "+d.FileMd5+" "+num)
          var log = '';
            for (var k in d.Lines) {
                var line = d.Lines[k];
                var content = line.replace(new RegExp('<','g'),"&lt;").replace(new RegExp('>','g'),"&gt;");
                if(content == ""){
                  content = "  "
                }
                log += '<dl class="'+get_class(line)+'">'
                    +'<dt>'+num+'</dt>'
                    +'<dd>'+content+'</dd>'
                    +'</dl>';
                num ++;
            }

            $(".file-log",$(".file_"+d.FileMd5)).append(log);
          continue;
        }
        var log = '<div class="one-file file_'+d.FileMd5+'">'
            +'<div class="title">'+d.Filename+'</div>'
            +'<div class="file-log">';
        var num = 1;
        for (var k in d.Lines) {
            var line = d.Lines[k];
            var content = line.replace(new RegExp('<','g'),"&lt;").replace(new RegExp('>','g'),"&gt;");
            if(content == ""){
                content = "  "
            }
            log += '<dl class="'+get_class(line)+'">'
                +'<dt>'+num+'</dt>'
                +'<dd>'+content+'</dd>'
                +'</dl>';
            num ++;
        }
        log +='</div></div>';
        str.push(log);
    }
    if(str.length > 0) {
        var html = str.join("");
        console.info(html.substr(0, 50));
        $(".table-responsive").append(html);
    }
    isLoad = false;
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