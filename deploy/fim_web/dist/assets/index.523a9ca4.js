import{d as k,b as V,p as A,l as M,a as p,c as h,f as e,h as g,g as a,w as _,q as j,s as I,t as $,v as F,e as m,r as N,x as U,E as C,y as L,i as x}from"./index.00362c69.js";import{M as R}from"./send_msg.d8fa9662.js";const T={class:"fim_slide"},S=["src"],q={class:"fim_menus"},P=k({__name:"fim_slide",setup(y){const s=V(),n=A();function u(c){n.push({name:c})}return(c,f)=>{const l=M("el-icon"),i=M("router-link");return p(),h("div",T,[e("div",{class:"avatar",onClick:f[0]||(f[0]=r=>u("info"))},[e("img",{src:g(s).userInfo.avatar,alt:""},null,8,S)]),e("div",q,[a(i,{class:"icon",to:{name:"session_welcome"}},{default:_(()=>[a(l,null,{default:_(()=>[a(g(j))]),_:1})]),_:1}),a(i,{class:"icon",to:{name:"contacts"}},{default:_(()=>[a(l,null,{default:_(()=>[a(g(I))]),_:1})]),_:1}),a(i,{class:"icon",to:{name:"notice"}},{default:_(()=>[a(l,null,{default:_(()=>[a(g($))]),_:1})]),_:1})]),a(i,{class:"other icon",to:{name:"info"}},{default:_(()=>[a(l,null,{default:_(()=>[a(g(F))]),_:1})]),_:1})])}}});const z={class:"video_action"},G={key:0,class:"info"},H={class:"name"},J=e("div",{class:"label"},"\u7B49\u5F85\u63A5\u542C",-1),K={class:"action"},Q=e("span",{class:"icon green"},[e("i",{class:"iconfont icon-shipin"})],-1),W=e("span",null,"\u63A5\u542C",-1),X=[Q,W],Y=e("span",{class:"icon red"},[e("i",{class:"iconfont icon-shipin"})],-1),Z=e("span",null,"\u6302\u65AD",-1),ee=[Y,Z],ae=k({__name:"fim_video_call_rev",setup(y){const s=V(),n=m(!1),u=m(),c=m(),f=m(),l=N({nickName:"",userID:0}),i=new R("user"),r=m(0);U(()=>s.chatMsgData,()=>{var o,t,v,D,b,w;if(s.chatMsgData.msg.type===7){if(((o=s.chatMsgData.msg.videoCallMsg)==null?void 0:o.flag)===2&&(l.userID=s.chatMsgData.sendUser.id,l.nickName=s.chatMsgData.sendUser.nickName,n.value=!0),!n.value)return;if(((t=s.chatMsgData.msg.videoCallMsg)==null?void 0:t.flag)===6){C.info("\u5BF9\u65B9\u7ED3\u675F\u89C6\u9891\u901A\u8BDD"),n.value=!1,r.value=0,c.value.srcObject=null,u.value.srcObject=null;return}if(((v=s.chatMsgData.msg.videoCallMsg)==null?void 0:v.flag)===3){C.info("\u53D1\u8D77\u8005\u5DF2\u6302\u65AD"),n.value=!1,r.value=0,c.value.srcObject=null,u.value.srcObject=null;return}switch((D=s.chatMsgData.msg.videoCallMsg)==null?void 0:D.type){case"offer":B((b=s.chatMsgData.msg.videoCallMsg)==null?void 0:b.data);break;case"answer_ice":d.value.addIceCandidate((w=s.chatMsgData.msg.videoCallMsg)==null?void 0:w.data);break}}},{deep:!0});const d=m();function B(o){var t;d.value.addTrack((t=f.value)==null?void 0:t.getVideoTracks()[0],f.value),d.value.setRemoteDescription(o).then(()=>{d.value.createAnswer().then(v=>{d.value.setLocalDescription(v).then(()=>{i.videoCallMsg(l.userID,{flag:9,type:"answer",data:v})})})})}function E(){navigator.mediaDevices.getUserMedia({audio:!0,video:!0}).then(o=>{u.value.srcObject=o,f.value=o,u.value.muted=!1,i.videoCallMsg(l.userID,{flag:3}),d.value=new RTCPeerConnection,d.value.addEventListener("icecandidate",t=>{t.candidate&&i.videoCallMsg(l.userID,{flag:9,type:"offer_ice",data:t.candidate})}),d.value.addEventListener("track",t=>{c.value.srcObject=t.streams[0],r.value=1,c.value.muted=!1})}).catch(o=>{C.error("\u65E0\u6743\u9650\u83B7\u53D6\u89C6\u9891\u4FE1\u606F")})}function O(){var o;(o=d.value)==null||o.close(),i.videoCallMsg(l.userID,{flag:r.value===1?4:2}),n.value=!1,r.value=0,c.value.srcObject=null,u.value.srcObject=null,C.info("\u7ED3\u675F\u89C6\u9891\u901A\u8BDD")}return(o,t)=>{const v=M("el-dialog");return p(),h("div",null,[a(v,{width:"400px","append-to-body":"",modelValue:n.value,"onUpdate:modelValue":t[0]||(t[0]=D=>n.value=D),"modal-class":"video_call_rev_dialog",draggable:""},{default:_(()=>[e("div",z,[r.value===0?(p(),h("div",G,[e("div",H,"\u6765\u81EA"+L(l.nickName)+"\u7684\u89C6\u9891\u901A\u8BDD",1),J])):x("",!0),e("div",K,[r.value===0?(p(),h("section",{key:0,onClick:E},X)):x("",!0),e("section",{onClick:O},ee)]),e("video",{class:"local_video",ref_key:"localVideo",ref:u,autoplay:"",muted:""},null,512),e("video",{class:"remote_video",ref_key:"remoteVideo",ref:c,autoplay:"",muted:""},null,512)])]),_:1},8,["modelValue"])])}}});const se={class:"fim_web"},te={class:"fim_main"},ne=k({__name:"index",setup(y){return(s,n)=>{const u=M("router-view");return p(),h("div",se,[a(P),a(ae),e("div",te,[a(u)])])}}});export{ne as default};
