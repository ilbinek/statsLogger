#include "kill.hpp"

Kill::Kill(std::string killerUID, std::string killedUID, std::string weapon, std::string distance, std::string time) {
    this->killerUID = killerUID;
    this->killedUID = killedUID;
    this->weapon = weapon;
    this->distance = distance;
    this->time = time;
};
