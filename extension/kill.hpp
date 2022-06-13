#ifndef KILL_H
#define KILL_H

#include <string>

class Kill {
    private:
        std::string killerUID;
        std::string killedUID;
        std::string weapon;
        std::string distance;
        std::string time;

    public:
        Kill(std::string killerUID, std::string killedUID, std::string weapon, std::string distance, std::string time);
        std::string getKillerUID() {return killerUID;};
        std::string getKilledUID() {return killedUID;};
        std::string getWeapon() {return weapon;};
        std::string getDistance() {return distance;};
        std::string getTime() {return time;};
};

#endif
