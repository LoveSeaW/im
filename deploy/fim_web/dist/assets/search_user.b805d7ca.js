import{d as E,u as F,b as A,r as i,x as N,e as U,c as r,f as l,F as V,z as I,g as c,i as S,l as n,a,w as u,m as _,y as j,K as f,T as M,E as O,k as P}from"./index.00362c69.js";import"./fim_avatar.vue_vue_type_style_index_0_lang.0abc880b.js";import"./search_dialog.vue_vue_type_style_index_0_lang.8710fe3f.js";import{s as T}from"./add_user_modal.edb5efa1.js";const $={class:"search_user_view"},K={class:"search_user_list"},L={class:"item"},R={class:"info"},q={class:"nickname"},G={key:0,class:"no_data"},Y=E({__name:"search_user",setup(H){F();const d=A(),s=i({key:"",limit:8,page:1}),o=i({list:[],count:0});async function p(){let e=await M(s);e.code&&O.error(e.msg),Object.assign(o,e.data)}N(()=>d.searchData,()=>{s.key=d.searchData.value,p()},{immediate:!0,deep:!0});function y(){p()}const k=U(!1),m=i({userID:0,nickname:"",abstract:"",avatar:"",isFriend:!1});function C(e){k.value=!0,Object.assign(m,e),T(m)}function x(e){P.push({name:"session_user_chat",params:{id:e.userID}})}return(e,g)=>{var v;const b=n("el-avatar"),w=n("el-text"),h=n("el-button"),z=n("el-empty"),B=n("el-pagination");return a(),r("div",$,[l("div",K,[(a(!0),r(V,null,I(o.list,t=>(a(),r("div",L,[c(b,{size:50,src:t.avatar},null,8,["src"]),l("div",R,[l("div",q,[c(w,{style:{width:"5rem"},truncated:""},{default:u(()=>[_(j(t.nickname),1)]),_:2},1024)]),t.isFriend?(a(),f(h,{key:0,onClick:D=>x(t),type:"primary",size:"small"},{default:u(()=>[_("\u53BB\u804A\u5929")]),_:2},1032,["onClick"])):(a(),f(h,{key:1,type:"primary",onClick:D=>C(t),size:"small"},{default:u(()=>[_("\u52A0\u597D\u53CB")]),_:2},1032,["onClick"]))])]))),256))]),((v=o.list)==null?void 0:v.length)===0?(a(),r("div",G,[c(z,{"image-size":200,description:"\u6682\u65E0\u6570\u636E"})])):S("",!0),c(B,{class:"search_page",onCurrentChange:y,"hide-on-single-page":"","current-page":s.page,"onUpdate:currentPage":g[0]||(g[0]=t=>s.page=t),"default-page-size":s.limit,layout:"prev, pager, next",total:o.count},null,8,["current-page","default-page-size","total"])])}}});export{Y as default};