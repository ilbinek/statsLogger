#include "player.hpp"

Player::Player(std::string uid, std::string name, std::string role, std::string className, std::string side, std::string squad) {
    this->uid = uid;
    this->name = name;
    this->role = role;
    this->className = className;
    this->side = side;
    this->squad = squad;
    this->shots = 0;
    this->hits = 0;
}
