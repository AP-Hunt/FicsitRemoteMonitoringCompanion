import Marker from "./marker"

export class MarkerSet<T extends Marker> extends Set<T> {

    public constructor(set : readonly T[] | MarkerSet<T> | null) {
        super(set)
    }

    static fromArray<T extends Marker>(elements : T[]): MarkerSet<T> {
        return new MarkerSet(elements);
    }

    isSuperset(subset : Set<T>) {
        for (let elem of subset) {
            if (!this.has(elem)) {
                return false
            }
        }
        return true
    }

    union(setB : Set<T>) {
        let _union : MarkerSet<T> = new MarkerSet(this)

        for (let elem of setB) {
            _union.add(elem)
        }
        return _union
    }

    intersection(setB : MarkerSet<T>) {
        let _intersection = new MarkerSet(this)
        for (let elem of setB) {
            if (this.has(elem)) {
                _intersection.add(elem)
            }
        }
        return _intersection
    }

    symmetricDifference(setB : MarkerSet<T>) {
        let _difference = new MarkerSet(this)
        for (let elem of setB) {
            if (_difference.has(elem)) {
                _difference.delete(elem)
            } else {
                _difference.add(elem)
            }
        }
        return _difference
    }

    difference(setB : MarkerSet<T>) {
        let _difference = new MarkerSet(this)
        for (let elem of setB) {
            _difference.delete(elem)
        }
        return _difference
    }
}