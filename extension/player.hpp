#ifndef PLAYER_H
#define PLAYER_H

#include <string>

class Player {
    private:
        std::string uid;
        std::string name;
        std::string role;
        std::string className;
        std::string side;
        std::string squad;
        int shots;
        int hits;

    public:
        Player(std::string uid, std::string name, std::string role, std::string className, std::string side, std::string squad);
        std::string getUID() {return uid;};
        std::string getName() {return name;};
        std::string getRole() {return role;};
        void setRole(std::string role) {this->role = role;};
        std::string getClassName() {return className;};
        void setClassName(std::string className) {this->className = className;};
        std::string getSide() {return side;};
        void setSide(std::string side) {this->side = side;};
        std::string getSquad() {return squad;};
        void setSquad(std::string squad) {this->squad = squad;};
        int getShots() {return shots;};
        void shot() {shots++;};
        int getHits() {return hits;};
        void hit() {hits++;};
};

#endif
