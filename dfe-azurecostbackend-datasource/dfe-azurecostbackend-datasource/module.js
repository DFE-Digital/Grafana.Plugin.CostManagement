define(["@grafana/data","@grafana/runtime","react","@grafana/ui"],((e,n,t,s)=>(()=>{"use strict";var a={305:n=>{n.exports=e},545:e=>{e.exports=n},388:e=>{e.exports=s},650:e=>{e.exports=t}},r={};function i(e){var n=r[e];if(void 0!==n)return n.exports;var t=r[e]={exports:{}};return a[e](t,t.exports,i),t.exports}i.n=e=>{var n=e&&e.__esModule?()=>e.default:()=>e;return i.d(n,{a:n}),n},i.d=(e,n)=>{for(var t in n)i.o(n,t)&&!i.o(e,t)&&Object.defineProperty(e,t,{enumerable:!0,get:n[t]})},i.o=(e,n)=>Object.prototype.hasOwnProperty.call(e,n),i.r=e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})};var c={};return(()=>{i.r(c),i.d(c,{plugin:()=>l});var e=i(305),n=i(545);const t={constant:6.5};class s extends n.DataSourceWithBackend{getDefaultQuery(e){return t}constructor(e){super(e)}}var a=i(650),r=i.n(a),o=i(388);const l=new e.DataSourcePlugin(s).setConfigEditor((function(e){const{onOptionsChange:n,options:t}=e,{secureJsonFields:s}=t,a=t.secureJsonData||{};return r().createElement("div",{className:"gf-form-group"},r().createElement(o.InlineField,{label:"Password / Client Secret",labelWidth:27},r().createElement(o.SecretInput,{isConfigured:s&&s.Password,value:a.Password||"",placeholder:"secure Password / Client Secret (backend only)",width:100,onReset:()=>{n(Object.assign(Object.assign({},t),{secureJsonFields:Object.assign(Object.assign({},t.secureJsonFields),{Password:!1}),secureJsonData:Object.assign(Object.assign({},t.secureJsonData),{Password:""})}))},onChange:e=>{n(Object.assign(Object.assign({},t),{secureJsonData:{Password:e.target.value}}))}})),r().createElement(o.InlineField,{label:"ClientID",labelWidth:12},r().createElement(o.SecretInput,{isConfigured:s&&s.ClientID,value:a.ClientID||"",placeholder:"secure Client ID (backend only)",width:100,onReset:()=>{n(Object.assign(Object.assign({},t),{secureJsonFields:Object.assign(Object.assign({},t.secureJsonFields),{ClientID:!1}),secureJsonData:Object.assign(Object.assign({},t.secureJsonData),{ClientID:""})}))},onChange:e=>{n(Object.assign(Object.assign({},t),{secureJsonData:{ClientID:e.target.value}}))}})),r().createElement(o.InlineField,{label:"TenantID",labelWidth:12},r().createElement(o.SecretInput,{isConfigured:s&&s.TenantID,value:a.TenantID||"",placeholder:"secure Tenant ID (backend only)",width:60,onReset:()=>{n(Object.assign(Object.assign({},t),{secureJsonFields:Object.assign(Object.assign({},t.secureJsonFields),{TenantID:!1}),secureJsonData:Object.assign(Object.assign({},t.secureJsonData),{TenantID:""})}))},onChange:e=>{n(Object.assign(Object.assign({},t),{secureJsonData:{TenantID:e.target.value}}))}})),r().createElement(o.InlineField,{label:"SubscriptionID",labelWidth:17},r().createElement(o.SecretInput,{isConfigured:s&&s.SubscriptionID,value:a.SubscriptionID||"",placeholder:"secure Subscription ID (backend only)",width:100,onReset:()=>{n(Object.assign(Object.assign({},t),{secureJsonFields:Object.assign(Object.assign({},t.secureJsonFields),{SubscriptionID:!1}),secureJsonData:Object.assign(Object.assign({},t.secureJsonData),{SubscriptionID:""})}))},onChange:e=>{n(Object.assign(Object.assign({},t),{secureJsonData:{SubscriptionID:e.target.value}}))}})),r().createElement(o.InlineField,{label:"Region",labelWidth:12},r().createElement(o.SecretInput,{isConfigured:s&&s.Region,value:a.Region||"",placeholder:"secure Region (backend only)",width:100,onReset:()=>{n(Object.assign(Object.assign({},t),{secureJsonFields:Object.assign(Object.assign({},t.secureJsonFields),{Region:!1}),secureJsonData:Object.assign(Object.assign({},t.secureJsonData),{Region:""})}))},onChange:e=>{n(Object.assign(Object.assign({},t),{secureJsonData:{Region:e.target.value}}))}})))})).setQueryEditor((function({query:e,onChange:n,onRunQuery:t}){const{queryText:s,constant:a}=e;return r().createElement("div",{className:"gf-form"},r().createElement(o.InlineField,{label:"Constant"},r().createElement(o.Input,{onChange:s=>{n(Object.assign(Object.assign({},e),{constant:parseFloat(s.target.value)})),t()},value:a,width:8,type:"number",step:"0.1"})),r().createElement(o.InlineField,{label:"Query Text",labelWidth:16,tooltip:"Not used yet"},r().createElement(o.Input,{onChange:t=>{n(Object.assign(Object.assign({},e),{queryText:t.target.value}))},value:s||""})))}))})(),c})()));
//# sourceMappingURL=module.js.map