import{d as E,b as L,e as u,x as z,c as _,g as i,w as d,F as V,z as N,aA as h,E as c,l as m,a as n,h as b,f as C,y as k,K as g,ap as T}from"./index.00362c69.js";import{_ as R}from"./avatar_cropper.vue_vue_type_style_index_0_lang.71cbc4ee.js";import"./file_api.3df89dab.js";const j={class:"my_info_view"},D={key:0},U=["onClick"],M=E({__name:"my_info",setup(W){const l=L(),a=u([{label:"\u6635\u79F0",isShowIpt:!1,maxlength:13,val:l.userConfInfo.nickname,type:"text",old:"",key:"nickname"},{label:"\u7B80\u4ECB",isShowIpt:!1,type:"textarea",rows:3,old:"",val:l.userConfInfo.abstract,key:"abstract"}]),p=u();z(()=>l.userConfInfo,()=>{a.value[0].val=l.userConfInfo.nickname,a.value[1].val=l.userConfInfo.abstract},{deep:!0});function I(e){a.value[e].isShowIpt=!0,a.value[e].old=a.value[e].val,T(()=>{p.value.length&&p.value[0].focus()})}async function F(e){if(a.value[e].isShowIpt=!1,a.value[e].old==a.value[e].val)return;let s={};s[a.value[e].key]=a.value[e].val;let r=await h(s);if(r.code){c.error(r.msg);return}c.success(a.value[e].label+"\u4FEE\u6539\u6210\u529F")}const o=u({type:"",allowTypeList:[],limitSize:1,fixedNumber:[],previewWidth:0}),y=u();function S(){o.value={type:"browserLogo",allowTypeList:["png","jpg","jpeg"],limitSize:1,fixedNumber:[1,1],previewWidth:100},y.value.uploadFile()}async function x(e){let s=await h({avatar:e});if(s.code){c.error(s.msg);return}l.userConfInfo.avatar=e,l.userInfo.avatar=e,l.saveToken(),c.success("\u7528\u6237\u5934\u50CF\u66F4\u65B0\u6210\u529F")}return(e,s)=>{const r=m("el-avatar"),f=m("el-form-item"),B=m("el-input");return n(),_("div",j,[i(R,{ref_key:"clipperRef",ref:y,type:o.value.type,"allow-type-list":o.value.allowTypeList,"limit-size":o.value.limitSize,"fixed-number":o.value.fixedNumber,"preview-width":o.value.previewWidth,onConfirm:x},null,8,["type","allow-type-list","limit-size","fixed-number","preview-width"]),i(f,{label:"\u5934\u50CF"},{default:d(()=>[i(r,{src:b(l).userConfInfo.avatar,onClick:S},null,8,["src"])]),_:1}),i(f,{label:"\u7528\u6237\u53F7"},{default:d(()=>[C("span",null,k(b(l).userConfInfo.userID),1)]),_:1}),(n(!0),_(V,null,N(a.value,(t,w)=>(n(),g(f,{label:t.label},{default:d(()=>[t.isShowIpt?(n(),g(B,{key:1,ref_for:!0,ref_key:"editRefList",ref:p,maxlength:t.maxlength,rows:t.rows,type:t.type,onBlur:v=>F(w),class:"edit_ipt",modelValue:t.val,"onUpdate:modelValue":v=>t.val=v,placeholder:"\u4FEE\u6539"+t.label},null,8,["maxlength","rows","type","onBlur","modelValue","onUpdate:modelValue","placeholder"])):(n(),_("span",D,k(t.val),1)),C("i",{class:"iconfont icon-bianji",onClick:v=>I(w)},null,8,U)]),_:2},1032,["label"]))),256))])}}});export{M as default};
