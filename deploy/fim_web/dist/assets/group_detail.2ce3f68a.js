import{d as y,b as x,u as B,r as F,x as I,c as E,f as e,g as a,w as s,y as n,k as u,l as i,a as l,m as r,K as p,i as h,M as D,E as _,N}from"./index.00362c69.js";import{_ as S}from"./fim_avatar.vue_vue_type_style_index_0_lang.0abc880b.js";const V={class:"group_detail_view"},M={class:"top"},O={class:"left"},T={class:"nickname"},G={class:"item"},j=e("span",{class:"label"},"ID",-1),q={class:"val"},z={class:"item"},K=e("span",{class:"label"},"\u4EBA\u6570",-1),L={class:"val"},P={class:"item"},R=e("span",{class:"label"},"\u7B80\u4ECB",-1),$={class:"val"},H={class:"right"},J={class:"more"},Q=e("i",{class:"iconfont icon-more"},null,-1),U={class:"avatar"},W={class:"bottom"},te=y({__name:"group_detail",setup(X){const f=x(),d=B(),t=F({groupId:0,title:"",abstract:"",memberCount:0,memberOnlineCount:0,avatar:"",creator:{userId:0,avatart:"",nickname:""},adminList:[],role:0,isProhibition:!1,prohibitionTime:void 0,isSearch:!1,isInvite:!1,isTemporarySession:!1});async function g(){let o=await D(Number(d.params.id));if(o.code){_.error(o.msg);return}Object.assign(t,o.data)}I(()=>d.params.id,()=>{g()},{immediate:!0});function m(){u.push({name:"session_group_chat",params:{id:t.groupId}})}function v(){u.push({name:"group_settings",params:{id:t.groupId}})}async function b(){let o=await N({memberId:f.userInfo.userID,id:t.groupId});if(o.code){_.error(o.msg);return}_.success("\u9000\u51FA\u7FA4\u804A\u6210\u529F"),u.push({name:"contacts"})}return(o,Y)=>{const w=i("el-text"),c=i("el-dropdown-item"),A=i("el-dropdown-menu"),C=i("el-dropdown"),k=i("el-button");return l(),E("div",V,[e("div",M,[e("div",O,[e("div",T,[a(w,{style:{"max-width":"30rem"},truncated:""},{default:s(()=>[r(n(t.title),1)]),_:1})]),e("div",G,[j,e("span",q,n(t.groupId),1)]),e("div",z,[K,e("span",L,n(t.memberOnlineCount)+"/"+n(t.memberCount),1)]),e("div",P,[R,e("span",$,n(t.abstract),1)])]),e("div",H,[e("div",J,[a(C,{trigger:"click"},{dropdown:s(()=>[a(A,null,{default:s(()=>[a(c,{onClick:m},{default:s(()=>[r("\u8FDB\u5165\u7FA4\u804A")]),_:1}),t.role!==3?(l(),p(c,{key:0,onClick:v},{default:s(()=>[r("\u7FA4\u8BBE\u7F6E")]),_:1})):h("",!0),t.role!==1?(l(),p(c,{key:1,onClick:b,style:{color:"red"}},{default:s(()=>[r("\u9000\u51FA\u7FA4\u804A")]),_:1})):h("",!0)]),_:1})]),default:s(()=>[Q]),_:1})]),e("div",U,[a(S,{src:t.avatar,shape:"square",size:60},null,8,["src"])])])]),e("div",W,[a(k,{type:"primary",onClick:m},{default:s(()=>[r("\u8FDB\u5165\u7FA4\u804A")]),_:1})])])}}});export{te as default};
