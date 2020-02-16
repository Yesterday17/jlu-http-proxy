package main

const ReauthHTML = `<!DOCTYPE html>
<html style="height: 100%;">
  <head>
    <meta charset="utf-8" />
    <title>登录失效</title>
    <meta name="renderer" content="webkit" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, maximum-scale=1"
    />
    <link
      rel="stylesheet"
      href="http://vpns.jlu.edu.cn/wengine-vpn/js/layui/css/layui.css"
      media="all"
    />

    <style>
      .flex {
        display: flex;
        width: 100%;
        height: 100%;
        background-color: #333;
      }

      .flex > div {
        margin: auto;
      }

      .hide {
        display: none !important;
      }
    </style>
  </head>
  <body class="flex">
    <div class="layui-card">
      <div class="layui-card-header">登录已失效，请重新登录。</div>
      <div class="layui-card-body">
        <button
          type="button"
          class="layui-btn layui-btn-fluid"
          onclick="reauth()"
        >
          <span
            id="loading"
            class="layui-icon layui-anim layui-anim-rotate layui-anim-loop layui-icon-loading hide"
          ></span>
          重新登录
        </button>
      </div>
    </div>

    <script
      src="http://vpns.jlu.edu.cn/wengine-vpn/js/layui/layui.js"
      charset="utf-8"
    ></script>
    <script>
      function reauth() {
        document.querySelector("#loading").classList.remove("hide");

        fetch("https://vpns.jlu.edu.cn/jlu-http-proxy/api/reauth")
          .then(d => d.json())
          .then(j => (j.success ? location.reload() : alert(j.message)));
      }
    </script>
  </body>
</html>`
