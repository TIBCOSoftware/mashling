package util

/**
Mashling gateway constants
*/
const Gateway_Definition_File_Name = "mashling.json"

const Gateway_Trigger_Config_Ref_Key = "config"
const Gateway_Trigger_Config_Prefix = "${configurations."
const Gateway_Trigger_Config_Suffix = "}"
const Gateway_Trigger_Handler_UseReplyHandler = "useReplyHandler"
const Gateway_Trigger_Handler_UseReplyHandler_Default = "false"
const Gateway_Trigger_Handler_AutoIdReply = "autoIdReply"
const Gateway_Trigger_Handler_AutoIdReply_Default = "false"
const Gateway_Trigger_Handler_If = "if"
const Gateway_Trigger_Metadata_JSON_Name = "trigger.json"
const Gateway_Trigger_Optimize_Property = "optimize"
const Gateway_Trigger_Optimize_Property_Default = false
const Gateway_Trigger_Setting_Env_Prefix = "${ENV."
const Gateway_Trigger_Setting_Env_Suffix = "}"

const Gateway_Link_Condition_LHS_Start_Expr = "${"
const Gateway_Link_Condition_LHS_End_Expr = "}"
const Gateway_JSON_Content_Root_Env_Key = "TRIGGER_CONTENT_ROOT"
const Gateway_Link_Condition_LHS_JSON_Content_Prefix_Default = "trigger.content"
const Gateway_Link_Condition_LHS_JSONPath_Root = "$"
const Gateway_Link_Condition_LHS_Header_Prifix = "trigger.header."
const Gateway_Link_Condition_LHS_Environment_Prifix = "env."

/**
Flogo constants
*/
const Flogo_App_Type = "flogo:app"
const Flogo_App_Embed_Config_Property = "FLOGO_EMBED"
const Flogo_App_Embed_Config_Property_Default = true
const Flogo_Trigger_Handler_Setting_Condition = "Condition"
