!function(e){var n=window.webpackHotUpdate;window.webpackHotUpdate=function(e,t){!function(e,n){if(!O[e]||!g[e])return;for(var t in g[e]=!1,n)Object.prototype.hasOwnProperty.call(n,t)&&(h[t]=n[t]);0==--m&&0===b&&I()}(e,t),n&&n(e,t)};var t,r=!0,o="ea8228deb34a4a796ec9",c={},i=[],a=[];function d(e){var n=x[e];if(!n)return k;var r=function(r){return n.hot.active?(x[r]?-1===x[r].parents.indexOf(e)&&x[r].parents.push(e):(i=[e],t=r),-1===n.children.indexOf(r)&&n.children.push(r)):(console.warn("[HMR] unexpected require("+r+") from disposed module "+e),i=[]),k(r)},o=function(e){return{configurable:!0,enumerable:!0,get:function(){return k[e]},set:function(n){k[e]=n}}};for(var c in k)Object.prototype.hasOwnProperty.call(k,c)&&"e"!==c&&"t"!==c&&Object.defineProperty(r,c,o(c));return r.e=function(e){return"ready"===u&&p("prepare"),b++,k.e(e).then(n,(function(e){throw n(),e}));function n(){b--,"prepare"===u&&(w[e]||j(e),0===b&&0===m&&I())}},r.t=function(e,n){return 1&n&&(e=r(e)),k.t(e,-2&n)},r}function s(n){var r={_acceptedDependencies:{},_declinedDependencies:{},_selfAccepted:!1,_selfDeclined:!1,_selfInvalidated:!1,_disposeHandlers:[],_main:t!==n,active:!0,accept:function(e,n){if(void 0===e)r._selfAccepted=!0;else if("function"==typeof e)r._selfAccepted=e;else if("object"==typeof e)for(var t=0;t<e.length;t++)r._acceptedDependencies[e[t]]=n||function(){};else r._acceptedDependencies[e]=n||function(){}},decline:function(e){if(void 0===e)r._selfDeclined=!0;else if("object"==typeof e)for(var n=0;n<e.length;n++)r._declinedDependencies[e[n]]=!0;else r._declinedDependencies[e]=!0},dispose:function(e){r._disposeHandlers.push(e)},addDisposeHandler:function(e){r._disposeHandlers.push(e)},removeDisposeHandler:function(e){var n=r._disposeHandlers.indexOf(e);n>=0&&r._disposeHandlers.splice(n,1)},invalidate:function(){switch(this._selfInvalidated=!0,u){case"idle":(h={})[n]=e[n],p("ready");break;case"ready":P(n);break;case"prepare":case"check":case"dispose":case"apply":(v=v||[]).push(n)}},check:E,apply:D,status:function(e){if(!e)return u;l.push(e)},addStatusHandler:function(e){l.push(e)},removeStatusHandler:function(e){var n=l.indexOf(e);n>=0&&l.splice(n,1)},data:c[n]};return t=void 0,r}var l=[],u="idle";function p(e){u=e;for(var n=0;n<l.length;n++)l[n].call(null,e)}var f,h,y,v,m=0,b=0,w={},g={},O={};function _(e){return+e+""===e?+e:e}function E(e){if("idle"!==u)throw new Error("check() is only allowed in idle status");return r=e,p("check"),(n=1e4,n=n||1e4,new Promise((function(e,t){if("undefined"==typeof XMLHttpRequest)return t(new Error("No browser support"));try{var r=new XMLHttpRequest,c=k.p+""+o+".hot-update.json";r.open("GET",c,!0),r.timeout=n,r.send(null)}catch(e){return t(e)}r.onreadystatechange=function(){if(4===r.readyState)if(0===r.status)t(new Error("Manifest request to "+c+" timed out."));else if(404===r.status)e();else if(200!==r.status&&304!==r.status)t(new Error("Manifest request to "+c+" failed."));else{try{var n=JSON.parse(r.responseText)}catch(e){return void t(e)}e(n)}}}))).then((function(e){if(!e)return p(H()?"ready":"idle"),null;g={},w={},O=e.c,y=e.h,p("prepare");var n=new Promise((function(e,n){f={resolve:e,reject:n}}));h={};return j(1),"prepare"===u&&0===b&&0===m&&I(),n}));var n}function j(e){O[e]?(g[e]=!0,m++,function(e){var n=document.createElement("script");n.charset="utf-8",n.src=k.p+""+e+"."+o+".hot-update.js",document.head.appendChild(n)}(e)):w[e]=!0}function I(){p("ready");var e=f;if(f=null,e)if(r)Promise.resolve().then((function(){return D(r)})).then((function(n){e.resolve(n)}),(function(n){e.reject(n)}));else{var n=[];for(var t in h)Object.prototype.hasOwnProperty.call(h,t)&&n.push(_(t));e.resolve(n)}}function D(n){if("ready"!==u)throw new Error("apply() is only allowed in ready status");return function n(r){var a,d,s,l,u;function f(e){for(var n=[e],t={},r=n.map((function(e){return{chain:[e],id:e}}));r.length>0;){var o=r.pop(),c=o.id,i=o.chain;if((l=x[c])&&(!l.hot._selfAccepted||l.hot._selfInvalidated)){if(l.hot._selfDeclined)return{type:"self-declined",chain:i,moduleId:c};if(l.hot._main)return{type:"unaccepted",chain:i,moduleId:c};for(var a=0;a<l.parents.length;a++){var d=l.parents[a],s=x[d];if(s){if(s.hot._declinedDependencies[c])return{type:"declined",chain:i.concat([d]),moduleId:c,parentId:d};-1===n.indexOf(d)&&(s.hot._acceptedDependencies[c]?(t[d]||(t[d]=[]),m(t[d],[c])):(delete t[d],n.push(d),r.push({chain:i.concat([d]),id:d})))}}}}return{type:"accepted",moduleId:e,outdatedModules:n,outdatedDependencies:t}}function m(e,n){for(var t=0;t<n.length;t++){var r=n[t];-1===e.indexOf(r)&&e.push(r)}}H();var b={},w=[],g={},E=function(){console.warn("[HMR] unexpected require("+I.moduleId+") to disposed module")};for(var j in h)if(Object.prototype.hasOwnProperty.call(h,j)){var I;u=_(j),I=h[j]?f(u):{type:"disposed",moduleId:j};var D=!1,P=!1,M=!1,A="";switch(I.chain&&(A="\nUpdate propagation: "+I.chain.join(" -> ")),I.type){case"self-declined":r.onDeclined&&r.onDeclined(I),r.ignoreDeclined||(D=new Error("Aborted because of self decline: "+I.moduleId+A));break;case"declined":r.onDeclined&&r.onDeclined(I),r.ignoreDeclined||(D=new Error("Aborted because of declined dependency: "+I.moduleId+" in "+I.parentId+A));break;case"unaccepted":r.onUnaccepted&&r.onUnaccepted(I),r.ignoreUnaccepted||(D=new Error("Aborted because "+u+" is not accepted"+A));break;case"accepted":r.onAccepted&&r.onAccepted(I),P=!0;break;case"disposed":r.onDisposed&&r.onDisposed(I),M=!0;break;default:throw new Error("Unexception type "+I.type)}if(D)return p("abort"),Promise.reject(D);if(P)for(u in g[u]=h[u],m(w,I.outdatedModules),I.outdatedDependencies)Object.prototype.hasOwnProperty.call(I.outdatedDependencies,u)&&(b[u]||(b[u]=[]),m(b[u],I.outdatedDependencies[u]));M&&(m(w,[I.moduleId]),g[u]=E)}var S,L=[];for(d=0;d<w.length;d++)u=w[d],x[u]&&x[u].hot._selfAccepted&&g[u]!==E&&!x[u].hot._selfInvalidated&&L.push({module:u,parents:x[u].parents.slice(),errorHandler:x[u].hot._selfAccepted});p("dispose"),Object.keys(O).forEach((function(e){!1===O[e]&&function(e){delete installedChunks[e]}(e)}));var T,U,q=w.slice();for(;q.length>0;)if(u=q.pop(),l=x[u]){var B={},R=l.hot._disposeHandlers;for(s=0;s<R.length;s++)(a=R[s])(B);for(c[u]=B,l.hot.active=!1,delete x[u],delete b[u],s=0;s<l.children.length;s++){var N=x[l.children[s]];N&&((S=N.parents.indexOf(u))>=0&&N.parents.splice(S,1))}}for(u in b)if(Object.prototype.hasOwnProperty.call(b,u)&&(l=x[u]))for(U=b[u],s=0;s<U.length;s++)T=U[s],(S=l.children.indexOf(T))>=0&&l.children.splice(S,1);p("apply"),void 0!==y&&(o=y,y=void 0);for(u in h=void 0,g)Object.prototype.hasOwnProperty.call(g,u)&&(e[u]=g[u]);var C=null;for(u in b)if(Object.prototype.hasOwnProperty.call(b,u)&&(l=x[u])){U=b[u];var W=[];for(d=0;d<U.length;d++)if(T=U[d],a=l.hot._acceptedDependencies[T]){if(-1!==W.indexOf(a))continue;W.push(a)}for(d=0;d<W.length;d++){a=W[d];try{a(U)}catch(e){r.onErrored&&r.onErrored({type:"accept-errored",moduleId:u,dependencyId:U[d],error:e}),r.ignoreErrored||C||(C=e)}}}for(d=0;d<L.length;d++){var X=L[d];u=X.module,i=X.parents,t=u;try{k(u)}catch(e){if("function"==typeof X.errorHandler)try{X.errorHandler(e)}catch(n){r.onErrored&&r.onErrored({type:"self-accept-error-handler-errored",moduleId:u,error:n,originalError:e}),r.ignoreErrored||C||(C=n),C||(C=e)}else r.onErrored&&r.onErrored({type:"self-accept-errored",moduleId:u,error:e}),r.ignoreErrored||C||(C=e)}}if(C)return p("fail"),Promise.reject(C);if(v)return n(r).then((function(e){return w.forEach((function(n){e.indexOf(n)<0&&e.push(n)})),e}));return p("idle"),new Promise((function(e){e(w)}))}(n=n||{})}function H(){if(v)return h||(h={}),v.forEach(P),v=void 0,!0}function P(n){Object.prototype.hasOwnProperty.call(h,n)||(h[n]=e[n])}var x={};function k(n){if(x[n])return x[n].exports;var t=x[n]={i:n,l:!1,exports:{},hot:s(n),parents:(a=i,i=[],a),children:[]};return e[n].call(t.exports,t,t.exports,d(n)),t.l=!0,t.exports}k.m=e,k.c=x,k.d=function(e,n,t){k.o(e,n)||Object.defineProperty(e,n,{enumerable:!0,get:t})},k.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},k.t=function(e,n){if(1&n&&(e=k(e)),8&n)return e;if(4&n&&"object"==typeof e&&e&&e.__esModule)return e;var t=Object.create(null);if(k.r(t),Object.defineProperty(t,"default",{enumerable:!0,value:e}),2&n&&"string"!=typeof e)for(var r in e)k.d(t,r,function(n){return e[n]}.bind(null,r));return t},k.n=function(e){var n=e&&e.__esModule?function(){return e.default}:function(){return e};return k.d(n,"a",n),n},k.o=function(e,n){return Object.prototype.hasOwnProperty.call(e,n)},k.p="",k.h=function(){return o},d(5)(k.s=5)}({5:function(e,n){var t,r=window.location.origin,o="localhost"===window.location.hostname?"ws":"wss",c="".concat(o,"://").concat(window.location.host,"/ws");window.addEventListener("load",(function(){var e=document.getElementById("lobby"),n=document.getElementById("progress-bar"),o=0;document.getElementById("lobby-message").innerHTML="Waiting for players to join",document.getElementById("unsuccessful-message").innerHTML="No available players at the moment.",document.getElementById("retry").setAttribute("href","/join"),setTimeout((function(){(t=new WebSocket("".concat(c,"/lobby"))).onopen=function(){},t.onclose=function(){t=null},t.onmessage=function(e){e.data&&("ping"===e.data?t.send("success"):window.location="".concat(r,"/g/").concat(e.data))},t.onerror=console.error}),1200);var i=setInterval((function(){(o+=1)>60&&(clearInterval(i),document.getElementById("unsuccessful").classList.remove("d-none"),e.classList.add("d-none"),t=null),n.setAttribute("aria-valuenow","".concat(o)),n.style.width="".concat(Math.round(o/.6),"%")}),1e3)}))}});
//# sourceMappingURL=lobby.js.map