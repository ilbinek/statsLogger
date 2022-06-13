
class CfgPatches {
    class STATSLOGGER {
        name = "StatsLogger";
        author = "ilbinek";
        authors[] = {"ilbinek"};
        url = "https://github.com/ilbinek";
        version = 0.1;
        versionStr = "0.1";
        versionAr[] = {0, 1};
        requiredAddons[] = {};
        requiredVersion = 2.04;
        units[] = {};
        weapons[] = {};
    };
};

class CfgFunctions {
    class STATSLOGGER {
        class null {
            file = "statslogger\functions";
            class init {preInit = 1;};
            class addEventMission {};
            class eh_connected {};
            class eh_killed {};
            class export {};
            class getEventWeaponText {};
            class eh_fired {};
            //class eh_hit {};
            class mission_end{};
        };
    };
};
