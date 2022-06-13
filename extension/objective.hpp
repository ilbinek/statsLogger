#ifndef OBJECTIVE_H
#define OBJECTIVE_H

#include <string>

class Objective {
    private:
        std::string name;
        std::string side;
        std::string points;
        std::string finished;

    public:
        Objective(std::string name, std::string side, std::string points, std::string finished) {
            this->name = name;
            this->side = side;
            this->points = points;
            this->finished = finished;
        };
        void setName(std::string name) {this->name = name;};
        void setSide(std::string side) {this->side = side;};
        void setPoints(std::string points) {this->points = points;};
        void setFinished(std::string finished) {this->finished = finished;};

        std::string getName() {return name;};
        std::string getSide() {return side;};
        std::string getPoints() {return points;};
        std::string getFinished() {return finished;};
};

#endif